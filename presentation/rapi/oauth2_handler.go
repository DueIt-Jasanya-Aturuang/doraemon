package rapi

import (
	"errors"
	"net/http"

	"github.com/jasanya-tech/jasanya-response-backend-golang/_error"
	"github.com/jasanya-tech/jasanya-response-backend-golang/response"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/presentation/rapi/helper"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/presentation/rapi/schema"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

func (p *Presenter) LoginWithGoogle(w http.ResponseWriter, r *http.Request) {
	req := new(schema.RequestLoginGoogle)
	err := helper.DecodeJson(r, req)
	if err != nil {
		helper.ErrorResponseEncode(w, err)
		return
	}

	err = req.Validation()
	if err != nil {
		helper.ErrorResponseEncode(w, err)
		return
	}

	userGoogle, err := p.oauth2Usecase.GoogleClaimUser(r.Context(), &usecase.RequestLoginWithGoogle{
		Token:  req.Token,
		Device: req.Device,
	})
	if err != nil {
		if errors.Is(err, usecase.InvalidTokenOauth) {
			err = _error.HttpErrString("invalid token", response.CM05)
		}
		helper.ErrorResponseEncode(w, err)
		return
	}

	appID := r.Header.Get(util.AppIDHeader)
	var auth *usecase.ResponseAuth

	if !userGoogle.ExistsUser {
		auth, err = p.authUsecase.Register(r.Context(), &usecase.RequestRegister{
			FullName:        userGoogle.Name,
			Username:        userGoogle.GivenName,
			Email:           userGoogle.Email,
			Password:        userGoogle.ID,
			RePassword:      userGoogle.ID,
			EmailVerifiedAt: true,
			AppID:           appID,
			Role:            1,
			RememberMe:      true,
		})
		if err != nil {
			if errors.Is(err, usecase.EmailIsExist) {
				err = _error.HttpErrMapOfSlices(map[string][]string{
					"email": {
						"email sudah terdaftar",
					},
				}, response.CM06)
			}
			if errors.Is(err, usecase.UsernameIsExist) {
				err = _error.HttpErrMapOfSlices(map[string][]string{
					"username": {
						"username sudah tersedia",
					},
				}, response.CM06)
			}
			helper.ErrorResponseEncode(w, err)
			return
		}

	} else {
		auth, err = p.authUsecase.Login(r.Context(), &usecase.RequestLogin{
			EmailOrUsername: userGoogle.Email,
			Password:        userGoogle.ID,
			RememberMe:      true,
			Oauth2:          true,
			AppID:           appID,
		})
		if err != nil {
			helper.ErrorResponseEncode(w, err)
			return
		}
	}

	profileResp := schema.ProfileDueitResponse{
		ProfileID: auth.Profile.ProfileID,
		Quote:     auth.Profile.Quote,
		Profesi:   auth.Profile.Profesi,
	}
	userResp := schema.UserResponse{
		ID:              auth.User.ID,
		FullName:        auth.User.FullName,
		Gender:          auth.User.Gender,
		Image:           auth.User.Image,
		Username:        auth.User.Username,
		Email:           auth.User.Email,
		EmailFormat:     auth.User.EmailFormat,
		PhoneNumber:     auth.User.PhoneNumber,
		EmailVerifiedAt: auth.User.EmailVerifiedAt,
	}
	tokenResp := schema.TokenResponse{
		Token: auth.Token.Token,
	}

	resp := map[string]any{
		"profile":         profileResp,
		"user_repository": userResp,
		"token":           tokenResp,
	}

	helper.SuccessResponseEncode(w, resp, "register successfully")
}
