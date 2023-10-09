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

type AuthHandlerImpl struct {
	authUsecase       usecase.AuthUsecase
	otpUsecase        usecase.OTPUsecase
	apiServiceUsecase usecase.ApiServiceUsecase
}

func NewAuthHandlerImpl(
	authUsecase usecase.AuthUsecase,
	otpUsecase usecase.OTPUsecase,
	apiServiceUsecase usecase.ApiServiceUsecase,
) *AuthHandlerImpl {
	return &AuthHandlerImpl{
		authUsecase:       authUsecase,
		otpUsecase:        otpUsecase,
		apiServiceUsecase: apiServiceUsecase,
	}
}

func (p *Presenter) Register(w http.ResponseWriter, r *http.Request) {
	req := new(schema.RequestRegister)
	err := helper.DecodeJson(r, &req)
	if err != nil {
		helper.ErrorResponseEncode(w, err)
		return
	}

	err = req.Validation()
	if err != nil {
		helper.ErrorResponseEncode(w, err)
		return
	}

	auth, err := p.authUsecase.Register(r.Context(), &usecase.RequestRegister{
		FullName:        req.FullName,
		Username:        req.Username,
		Email:           req.Email,
		Password:        req.Password,
		RePassword:      req.RePassword,
		EmailVerifiedAt: false,
		AppID:           r.Header.Get(util.AppIDHeader),
		Role:            1,
		RememberMe:      false,
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

	_ = p.otpUsecase.Generate(r.Context(), &usecase.RequestGenerateOTP{
		Email:  req.Email,
		Type:   util.ActivasiAccount,
		UserID: auth.User.ID,
	})

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
		"profile": profileResp,
		"user":    userResp,
		"token":   tokenResp,
	}

	helper.SuccessResponseEncode(w, resp, "register successfully")
}

func (p *Presenter) Login(w http.ResponseWriter, r *http.Request) {
	req := new(schema.RequestLogin)
	err := helper.DecodeJson(r, &req)
	if err != nil {
		helper.ErrorResponseEncode(w, err)
		return
	}

	appID := r.Header.Get(util.AppIDHeader)

	err = req.Validation()
	if err != nil {
		helper.ErrorResponseEncode(w, err)
		return
	}

	auth, err := p.authUsecase.Login(r.Context(), &usecase.RequestLogin{
		EmailOrUsername: req.EmailOrUsername,
		Password:        req.Password,
		RememberMe:      req.RememberMe,
		Oauth2:          false,
		AppID:           appID,
	})
	if err != nil {
		if errors.Is(err, usecase.InvalidEmailOrUsernameOrPassword) {
			err = _error.HttpErrMapOfSlices(map[string][]string{
				"email_or_username": {
					"invalid email atau password",
				},
				"password": {
					"invalid email atau password",
				},
			}, response.CM06)
		}
		helper.ErrorResponseEncode(w, err)
		return
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
		"profile": profileResp,
		"user":    userResp,
		"token":   tokenResp,
	}

	helper.SuccessResponseEncode(w, resp, "login successfully")
}
