package rapi

import (
	"errors"
	"net/http"

	"github.com/jasanya-tech/jasanya-response-backend-golang/_error"
	"github.com/jasanya-tech/jasanya-response-backend-golang/response"
	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/presentation/rapi/helper"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/presentation/rapi/middleware"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/presentation/rapi/schema"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

func (p *Presenter) GenerateOTP(w http.ResponseWriter, r *http.Request) {
	req := new(schema.RequestGenerateOTP)
	err := helper.DecodeJson(r, req)
	if err != nil {
		helper.ErrorResponseEncode(w, err)
		return
	}

	// validasi request
	err = req.Validation()
	if err != nil {
		helper.ErrorResponseEncode(w, err)
		return
	}

	typeHeader := r.Header.Get(util.TypeHeader)
	userIDHeader := r.Header.Get(util.UserIDHeader)

	if typeHeader != util.ActivasiAccount && typeHeader != util.ForgotPassword {
		log.Info().Msgf("invalid type header | %s", typeHeader)
		helper.SuccessResponseEncode(w, nil, "kode otp telah berhasil dikirim, silahkan cek gmail anda")
		return
	}
	if typeHeader == util.ActivasiAccount {
		if err = util.ParseUUID(userIDHeader); err != nil {
			helper.ErrorResponseEncode(w, _error.HttpErrString("invalid user id", response.CM05))
			return
		}
	}

	err = p.otpUsecase.Generate(r.Context(), &usecase.RequestGenerateOTP{
		Email:  req.Email,
		Type:   typeHeader,
		UserID: userIDHeader,
	})
	if err != nil {
		if errors.Is(err, usecase.InvalidUserID) {
			err = _error.HttpErrString("invalid user id", response.CM05)
		}
		if errors.Is(err, usecase.InvalidEmail) {
			err = _error.HttpErrMapOfSlices(map[string][]string{
				"email": {
					"email tidak terdaftar",
				},
			}, response.CM06)
		}
		if errors.Is(err, usecase.EmailIsActivited) {
			err = _error.HttpErrMapOfSlices(map[string][]string{
				"email": {
					"permintaan anda tidak dapat di proses, email anda sudah di aktivasi silahkan login",
				},
			}, response.CM06)
		}
		helper.ErrorResponseEncode(w, err)
		return
	}

	// set limiter
	err = middleware.RateLimiterOTP(&usecase.RequestGenerateOTP{
		Email:  req.Email,
		Type:   typeHeader,
		UserID: userIDHeader,
	})
	if err != nil {
		helper.ErrorResponseEncode(w, err)
		return
	}

	helper.SuccessResponseEncode(w, nil, "kode otp telah berhasil dikirim, silahkan cek gmail anda")
}
