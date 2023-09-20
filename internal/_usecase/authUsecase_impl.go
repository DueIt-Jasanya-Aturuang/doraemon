package _usecase

import (
	"context"
	"database/sql"
	"errors"

	uuid "github.com/satori/go.uuid"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/converter"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/helper"
)

type AuthUsecaseImpl struct {
	userRepo     domain.UserRepository
	accessRepo   domain.AccessRepository
	apiService   domain.MicroServiceRepository
	securityRepo domain.SecurityRepository
}

func NewAuthUsecaseImpl(
	userRepo domain.UserRepository,
	accessRepo domain.AccessRepository,
	apiService domain.MicroServiceRepository,
	securityRepo domain.SecurityRepository,
) domain.AuthUsecase {
	return &AuthUsecaseImpl{
		userRepo:     userRepo,
		accessRepo:   accessRepo,
		apiService:   apiService,
		securityRepo: securityRepo,
	}
}

func (a *AuthUsecaseImpl) Login(ctx context.Context, req *domain.RequestLogin) (*domain.ResponseAuth, error) {
	if err := a.userRepo.OpenConn(ctx); err != nil {
		return nil, err
	}
	defer a.userRepo.CloseConn()

	user, err := a.userRepo.GetByEmailOrUsername(ctx, req.EmailOrUsername)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, InvalidEmailOrUsernameOrPassword
		}
		return nil, err
	}

	if !req.Oauth2 {
		checkPassword := helper.BcryptPasswordCompare(req.Password, user.Password)
		if !checkPassword {
			return nil, InvalidEmailOrUsernameOrPassword
		}
	}

	profile, err := a.apiService.GetProfileByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	rtat, err := helper.GenerateRTAT(user.ID, req.AppID, req.RememberMe)
	if err != nil {
		return nil, err
	}

	err = a.userRepo.StartTx(ctx, helper.LevelReadCommitted(), func() error {
		err = a.securityRepo.Create(ctx, rtat)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	emailFormat := helper.EmailFormat(user.Email)
	resp := converter.AuthModelToResp(user, profile, emailFormat)
	resp.Token = rtat.AcceesToken

	return resp, nil
}

func (a *AuthUsecaseImpl) Register(ctx context.Context, req *domain.RequestRegister) (*domain.ResponseAuth, error) {
	if err := a.userRepo.OpenConn(ctx); err != nil {
		return nil, err
	}
	defer a.userRepo.CloseConn()

	domain.CheckUserByEmail = req.Email
	exist, err := a.userRepo.Check(ctx, domain.CheckUserByEmail)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, EmailIsExist
	}

	domain.CheckUserByUsername = req.Username
	exist, err = a.userRepo.Check(ctx, domain.CheckUserByUsername)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, UsernameIsExist
	}

	passwordHash, err := helper.BcryptPasswordHash(req.Password)
	if err != nil {
		return nil, err
	}
	req.Password = passwordHash

	id := uuid.NewV4().String()
	user, access := converter.RegisterReqToModel(req, id)
	var profile *domain.Profile
	var at string

	err = a.userRepo.StartTx(ctx, helper.LevelReadCommitted(), func() error {
		err = a.userRepo.Create(ctx, user)
		if err != nil {
			return err
		}

		err = a.accessRepo.Create(ctx, access)
		if err != nil {
			return err
		}

		rtat, err := helper.GenerateRTAT(user.ID, req.AppID, false)
		if err != nil {
			return err
		}
		at = rtat.AcceesToken

		err = a.securityRepo.Create(ctx, rtat)
		if err != nil {
			return err
		}

		profileByte, err := helper.SerializeProfile(user.ID)
		if err != nil {
			return err
		}

		profile, err = a.apiService.CreateProfile(profileByte)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	emailFormat := helper.EmailFormat(user.Email)
	resp := converter.AuthModelToResp(user, profile, emailFormat)
	resp.Token = at

	return resp, nil
}
