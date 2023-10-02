package rest

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jasanya-tech/jasanya-response-backend-golang/_error"
	"github.com/jasanya-tech/jasanya-response-backend-golang/response"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/api/rest/helper"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/api/rest/middleware"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/api/validation"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

type OTPHandlerImpl struct {
	otpUsecase domain.OTPUsecase
}

func NewOTPHandlerImpl(
	otpUsecase domain.OTPUsecase,
) *OTPHandlerImpl {
	return &OTPHandlerImpl{
		otpUsecase: otpUsecase,
	}
}

func (h *OTPHandlerImpl) GenerateOTP(w http.ResponseWriter, r *http.Request) {
	req := new(domain.RequestGenerateOTP)
	err := helper.DecodeJson(r, req)
	if err != nil {
		helper.ErrorResponseEncode(w, err)
		return
	}
	req.Type = r.Header.Get(util.TypeHeader)

	// validasi request
	err = validation.GenerateOTPValidation(req)
	if err != nil {
		helper.ErrorResponseEncode(w, err)
		return
	}

	if req.Type == util.ActivasiAccount {
		userID := r.Header.Get("User-ID")
		if _, err := uuid.Parse(userID); err != nil {
			helper.ErrorResponseEncode(w, _error.HttpErrString("invalid user id", response.CM05))
			return
		}

		req.UserID = userID
	}

	// generate otp
	err = h.otpUsecase.Generate(r.Context(), req)
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
	err = middleware.RateLimiterOTP(req)
	if err != nil {
		helper.ErrorResponseEncode(w, err)
		return
	}

	helper.SuccessResponseEncode(w, nil, "kode otp telah berhasil dikirim, silahkan cek gmail anda")
}
