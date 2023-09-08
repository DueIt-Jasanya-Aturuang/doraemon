package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/helper"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/helper/conv"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util"
	_error "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/error"
)

type AuthUsecaseImpl struct {
	userRepo   repository.UserSqlRepo
	accessRepo repository.AccessSqlRepo
	accountApi repository.AccountApiRepo
}

func NewAuthUsecaseImpl(
	userRepo repository.UserSqlRepo,
	accessRepo repository.AccessSqlRepo,
	accountApi repository.AccountApiRepo,
) usecase.AuthUsecase {
	return &AuthUsecaseImpl{
		userRepo:   userRepo,
		accessRepo: accessRepo,
		accountApi: accountApi,
	}
}

func (a *AuthUsecaseImpl) Login(ctx context.Context, req *dto.LoginReq) (userResp *dto.UserResp, profileResp *dto.ProfileResp, err error) {
	err = a.userRepo.OpenConn(ctx)
	if err != nil {
		return nil, nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer a.userRepo.CloseConn()

	user, err := a.userRepo.GetUserByEmailOrUsername(ctx, req.EmailOrUsername)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Info().Msgf("user login tapi email or username nya tidak tersedia")
			return nil, nil, _error.BadLogin()
		}
		return nil, nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	if !req.Oauth2 {
		checkPassword := helper.BcryptPasswordCompare(req.Password, user.Password)
		if !checkPassword {
			log.Info().Msgf("password user tidak sama dengan yang ada di database")
			return nil, nil, _error.BadLogin()
		}
	}

	profile, err := a.accountApi.GetProfileByUserID(user.ID)
	if err != nil {
		return nil, nil, _error.ErrStringDefault(http.StatusBadGateway)
	}

	emailFormat := util.EmailFormat(user.Email)
	userResp, profileResp = conv.LoginModelToResp(user, profile, emailFormat)

	return userResp, profileResp, nil
}

func (a *AuthUsecaseImpl) Register(ctx context.Context, req *dto.RegisterReq) (userResp *dto.UserResp, err error) {
	err = a.userRepo.OpenConn(ctx)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer a.userRepo.CloseConn()

	exists, err := a.userRepo.CheckUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}
	if exists {
		log.Info().Msgf("user register tapi email sudah terdaftar")
		return nil, _error.BadExistField("email", "email has been registered")
	}

	exists, err = a.userRepo.CheckUserByUsername(ctx, req.Email)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}
	if exists {
		log.Info().Msgf("user register tapi username sudah terdaftar")
		return nil, _error.BadExistField("username", "username has been registered")
	}

	passwordHash, err := helper.BcryptPasswordHash(req.Password)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}
	req.Password = passwordHash

	err = a.userRepo.StartTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer func() {
		errEndTx := a.userRepo.EndTx(err)
		if errEndTx != nil {
			err = _error.ErrStringDefault(http.StatusInternalServerError)
			userResp = nil
		}
	}()

	id := uuid.NewV4().String()
	userConv, accessConv := conv.RegisterReqToModel(req, id)

	err = a.userRepo.CreateUser(ctx, userConv)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	_, err = a.accessRepo.CreateAccess(ctx, accessConv)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	profileReq := dto.ProfileReq{
		UserID: userConv.ID,
	}
	profileJson, err := json.Marshal(profileReq)
	if err != nil {
		log.Err(err).Msg("failed marshal profile req")
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	_, err = a.accountApi.CreateProfile(profileJson)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusBadGateway)
	}

	emailFormat := util.EmailFormat(userConv.Email)
	userResp = conv.RegisterModelToResp(userConv, emailFormat)

	return userResp, nil
}
