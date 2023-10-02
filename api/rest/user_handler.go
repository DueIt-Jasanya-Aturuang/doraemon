package rest

import (
	"errors"
	"net/http"

	"github.com/jasanya-tech/jasanya-response-backend-golang/_error"
	"github.com/jasanya-tech/jasanya-response-backend-golang/response"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/api/rest/helper"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/api/validation"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

type UserHandlerImpl struct {
	userUsecase domain.UserUsecase
	otpUsecase  domain.OTPUsecase
}

func NewUserHandlerImpl(
	userUsecase domain.UserUsecase,
	otpUsecase domain.OTPUsecase,
) *UserHandlerImpl {
	return &UserHandlerImpl{
		userUsecase: userUsecase,
		otpUsecase:  otpUsecase,
	}
}

func (h *UserHandlerImpl) ChangePassword(w http.ResponseWriter, r *http.Request) {
	req := new(domain.RequestChangePassword)
	err := helper.DecodeJson(r, req)
	if err != nil {
		helper.ErrorResponseEncode(w, err)
		return
	}

	userID := r.Header.Get(util.UserIDHeader)
	req.UserID = userID

	err = validation.ChangePasswordValidation(req)
	if err != nil {
		helper.ErrorResponseEncode(w, err)
		return
	}

	// reset password process
	err = h.userUsecase.ChangePassword(r.Context(), req)
	if err != nil {
		if errors.Is(err, usecase.InvalidUserID) {
			err = _error.HttpErrString("user id tidak valid", response.CM04)
		}
		if errors.Is(err, usecase.InvalidOldPassword) {
			err = _error.HttpErrMapOfSlices(map[string][]string{
				"old_password": {
					"password lama tidak sesuai",
				},
			}, response.CM06)

		}
		helper.ErrorResponseEncode(w, err)
		return
	}

	helper.SuccessResponseEncode(w, nil, "password anda telah berhasil dirubah")

}

func (h *UserHandlerImpl) ChangeUsername(w http.ResponseWriter, r *http.Request) {
	req := new(domain.RequestChangeUsername)
	err := helper.DecodeJson(r, req)
	if err != nil {
		helper.ErrorResponseEncode(w, err)
		return
	}

	userID := r.Header.Get(util.UserIDHeader)
	req.UserID = userID

	err = validation.ChangeUsernameValidation(req)
	if err != nil {
		helper.ErrorResponseEncode(w, err)
		return
	}

	err = h.userUsecase.ChangeUsername(r.Context(), req)
	if err != nil {
		if errors.Is(err, usecase.InvalidUserID) {
			err = _error.HttpErrString("user id tidak valid", response.CM04)
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

	helper.SuccessResponseEncode(w, map[string]string{
		"username": req.Username,
	}, "username anda telah berhasil dirubah")

}

func (h *UserHandlerImpl) ForgottenPassword(w http.ResponseWriter, r *http.Request) {
	req := new(domain.RequestValidationOTP)
	err := helper.DecodeJson(r, req)
	if err != nil {
		helper.ErrorResponseEncode(w, err)
		return
	}
	req.Type = util.ForgotPassword

	// validasi request
	err = validation.OTPValidation(req)
	if err != nil {
		helper.ErrorResponseEncode(w, err)
		return
	}

	// process validasi otp
	err = h.otpUsecase.Validation(r.Context(), req)
	if err != nil {
		if errors.Is(err, usecase.InvalidEmailOrOTP) {
			err = _error.HttpErrMapOfSlices(map[string][]string{
				"email": {
					"invalid email or otp",
				},
				"otp": {
					"invalid email or otp",
				},
			}, response.CM06)
		}
		helper.ErrorResponseEncode(w, err)
		return
	}

	// declare request forgotten password
	reqFP := new(domain.RequestForgottenPassword)
	reqFP.Email = req.Email

	// process forgotten password
	url, err := h.userUsecase.ForgottenPassword(r.Context(), reqFP)
	if err != nil {
		if errors.Is(err, usecase.InvalidUserID) {
			err = _error.HttpErrString("user id tidak valid", response.CM04)
		}
		helper.ErrorResponseEncode(w, err)
		return
	}

	data := map[string]string{
		"url_forgot_password": url,
	}

	helper.SuccessResponseEncode(w, data, "link forgot password")
}

func (h *UserHandlerImpl) ResetForgottenPassword(w http.ResponseWriter, r *http.Request) {
	req := new(domain.RequestResetForgottenPassword)
	err := helper.DecodeJson(r, req)
	if err != nil {
		helper.ErrorResponseEncode(w, err)
		return
	}

	// declare variable from url query param
	email := r.URL.Query().Get("email")
	token := r.URL.Query().Get("token")
	req.Email = email
	req.Token = token

	// validasi request
	err = validation.ResetForgottenPasswordValidation(req)
	if err != nil {
		helper.ErrorResponseEncode(w, err)
		return
	}

	// process reset forgotten password
	err = h.userUsecase.ResetForgottenPassword(r.Context(), req)
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

func (h *UserHandlerImpl) ActivasiAccount(w http.ResponseWriter, r *http.Request) {
	req := new(domain.RequestValidationOTP)
	err := helper.DecodeJson(r, req)
	if err != nil {
		helper.ErrorResponseEncode(w, err)
		return
	}
	req.Type = util.ActivasiAccount

	// validasi request otp
	err = validation.OTPValidation(req)
	if err != nil {
		helper.ErrorResponseEncode(w, err)
		return
	}

	// validasi otp
	err = h.otpUsecase.Validation(r.Context(), req)
	if err != nil {
		if errors.Is(err, usecase.InvalidEmailOrOTP) {
			err = _error.HttpErrMapOfSlices(map[string][]string{
				"email": {
					"invalid email or otp",
				},
				"otp": {
					"invalid email or otp",
				},
			}, response.CM06)
		}
		helper.ErrorResponseEncode(w, err)
		return
	}

	// process activasi account
	activasi, err := h.userUsecase.ActivasiAccount(r.Context(), req.Email)
	if err != nil {
		if errors.Is(err, usecase.InvalidUserID) {
			err = _error.HttpErrString("user id tidak valid", response.CM04)
		}
		if errors.Is(err, usecase.EmailIsActivited) {
			err = _error.HttpErrMapOfSlices(map[string][]string{
				"email": {
					"email anda sudah aktif silangkah login",
				},
			}, response.CM06)
		}
		helper.ErrorResponseEncode(w, err)
		return
	}

	helper.SuccessResponseEncode(w, activasi, "activasi berhasil")
}
