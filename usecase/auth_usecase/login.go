package auth_usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
)

func (a *AuthUsecaseImpl) Login(ctx context.Context, req *usecase.RequestLogin) (*usecase.ResponseAuth, error) {
	user, err := a.userRepo.GetByEmailOrUsername(ctx, req.EmailOrUsername)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, usecase.InvalidEmailOrUsernameOrPassword
		}
		return nil, err
	}

	if !req.Oauth2 {
		checkPassword := usecase.BcryptPasswordCompare(req.Password, user.Password)
		if !checkPassword {
			return nil, usecase.InvalidEmailOrUsernameOrPassword
		}
	}

	if err != nil {
		return nil, err
	}

	profile, err := a.apiServiceUsecase.GetProfileByUserID(user.ID, req.AppID)
	if err != nil {
		return nil, err
	}

	securityToken, err := a.securityUsecase.GenerateJWT(ctx, &usecase.RequestGenerateJWT{
		AppID:      req.AppID,
		UserID:     user.ID,
		RememberMe: req.RememberMe,
	})

	emailFormat := usecase.FormatEmail(user.Email)

	resp := usecase.AuthMergeModelToResponse{
		User:        user,
		Profile:     profile,
		FormatEmail: emailFormat,
		Token:       securityToken.Token,
	}

	return resp.Execute(), nil
}
