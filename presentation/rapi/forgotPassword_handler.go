package rapi

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/jasanya-tech/jasanya-response-backend-golang/_error"
	"github.com/jasanya-tech/jasanya-response-backend-golang/response"
	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/presentation/rapi/helper"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/presentation/rapi/schema"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

func (p *Presenter) ForgottenPassword(w http.ResponseWriter, r *http.Request) {
	req := new(schema.RequestValidationOTP)
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

	err = p.otpUsecase.Validation(r.Context(), &usecase.RequestValidationOTP{
		Email: req.Email,
		OTP:   req.OTP,
		Type:  util.ForgotPassword,
	})
	if err != nil {
		if errors.Is(err, usecase.InvalidEmailOrOTP) {
			err = _error.HttpErrMapOfSlices(map[string][]string{
				"email": {
					"invalid email or otp_usecase",
				},
				"otp_usecase": {
					"invalid email or otp_usecase",
				},
			}, response.CM06)
		}
		helper.ErrorResponseEncode(w, err)
		return
	}

	url, err := p.userUsecase.ForgottenPassword(r.Context(), &usecase.RequestForgottenPassword{
		Email: req.Email,
	})
	if err != nil {
		if errors.Is(err, usecase.InvalidUserID) {
			err = _error.HttpErrString("user_repository id tidak valid", response.CM04)
		}
		helper.ErrorResponseEncode(w, err)
		return
	}

	data := map[string]string{
		"url_forgot_password": url,
	}

	helper.SuccessResponseEncode(w, data, "link forgot password")
}

func (p *Presenter) ResetForgottenPassword(w http.ResponseWriter, r *http.Request) {
	req := new(schema.RequestResetForgottenPassword)
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

	email := r.URL.Query().Get("email")
	token := r.URL.Query().Get("token")

	match, err := regexp.MatchString(`^([A-Za-z.]|[0-9])+@gmail.com$`, email)
	if err != nil || !match {
		log.Info().Msgf("invalid email | email : ", email)
		helper.ErrorResponseEncode(w, _error.HttpErrString("invalid email in query param", response.CM05))
		return
	}
	if token == "" {
		log.Info().Msgf("invalid token | token : ", token)
		helper.ErrorResponseEncode(w, _error.HttpErrString("invalid token in query param", response.CM05))
		return
	}

	err = p.userUsecase.ResetForgottenPassword(r.Context(), &usecase.RequestResetForgottenPassword{
		Email:      email,
		Token:      token,
		Password:   req.Password,
		RePassword: req.RePassword,
	})
	if err != nil {
		if errors.Is(err, usecase.InvalidToken) {
			err = _error.HttpErrString("invalid token", response.CM04)
		}
		if errors.Is(err, usecase.TokenExpired) {
			err = _error.HttpErrString("token anda sudah expired", response.CM05)
		}
		helper.ErrorResponseEncode(w, err)
		return
	}

	helper.SuccessResponseEncode(w, nil, "password anda telah berhasil diubah")
}
