package auth_usecase

import (
	"context"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
)

func (a *AuthUsecaseImpl) Register(ctx context.Context, req *usecase.RequestRegister) (*usecase.ResponseAuth, error) {
	repository.CheckUserByEmail = req.Email
	exist, err := a.userRepo.Check(ctx, repository.CheckUserByEmail)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, usecase.EmailIsExist
	}

	repository.CheckUserByUsername = req.Username
	exist, err = a.userRepo.Check(ctx, repository.CheckUserByUsername)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, usecase.UsernameIsExist
	}

	passwordHash, err := usecase.BcryptPasswordHash(req.Password)
	if err != nil {
		return nil, err
	}
	req.Password = passwordHash

	user, access := usecase.AuthRegisterRequestToModel(req)
	var securityToken *usecase.ResponseJWT
	var profile *usecase.ResponseProfileDueit

	err = a.userRepo.StartTx(ctx, repository.LevelReadCommitted(), func() error {
		profile, err = a.apiServiceUsecase.GetProfileByUserID(user.ID, req.AppID)
		if err != nil {
			return err
		}

		err = a.userRepo.Create(ctx, user)
		if err != nil {
			return err
		}

		err = a.accessRepo.Create(ctx, access)
		if err != nil {
			return err
		}

		securityToken, err = a.securityUsecase.GenerateJWT(ctx, &usecase.RequestGenerateJWT{
			AppID:      req.AppID,
			UserID:     user.ID,
			RememberMe: req.RememberMe,
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	emailFormat := usecase.FormatEmail(user.Email)
	resp := usecase.AuthMergeModelToResponse{
		User:        user,
		Profile:     profile,
		FormatEmail: emailFormat,
		Token:       securityToken.Token,
	}

	return resp.Execute(), nil
}
