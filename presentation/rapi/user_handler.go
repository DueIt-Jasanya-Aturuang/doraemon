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

func (p *Presenter) ChangePassword(w http.ResponseWriter, r *http.Request) {
	req := new(schema.RequestChangePassword)
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

	userID := r.Header.Get(util.UserIDHeader)
	err = util.ParseUUID(userID)
	if err != nil {
		helper.ErrorResponseEncode(w, _error.HttpErrString("invalid user id", response.CM05))
		return
	}

	// reset password process
	err = p.userUsecase.ChangePassword(r.Context(), &usecase.RequestChangePassword{
		OldPassword: req.OldPassword,
		Password:    req.Password,
		RePassword:  req.RePassword,
		UserID:      userID,
	})
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

func (p *Presenter) ChangeUsername(w http.ResponseWriter, r *http.Request) {
	req := new(schema.RequestChangeUsername)
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

	userID := r.Header.Get(util.UserIDHeader)
	if err = util.ParseUUID(userID); err != nil {
		helper.ErrorResponseEncode(w, _error.HttpErrString("invalid user id", response.CM05))
		return
	}

	err = p.userUsecase.ChangeUsername(r.Context(), &usecase.RequestChangeUsername{
		Username: req.Username,
		UserID:   userID,
	})
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

func (p *Presenter) ActivasiAccount(w http.ResponseWriter, r *http.Request) {
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

	typeHeader := util.ActivasiAccount
	if err = util.TypeHeaderValidation(typeHeader); err != nil {
		helper.ErrorResponseEncode(w, err)
		return
	}

	err = p.otpUsecase.Validation(r.Context(), &usecase.RequestValidationOTP{
		Email: req.Email,
		OTP:   req.OTP,
		Type:  typeHeader,
	})
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

	activasi, err := p.userUsecase.ActivasiAccount(r.Context(), req.Email)
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

	resp := map[string]bool{
		"activited": activasi,
	}
	helper.SuccessResponseEncode(w, resp, "activasi berhasil")
}

func (p *Presenter) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("User-ID")

	if err := util.ParseUUID(userID); err != nil {
		helper.ErrorResponseEncode(w, _error.HttpErrString("invalid user id", response.CM05))
		return
	}

	user, err := p.userUsecase.GetUserByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, usecase.InvalidUserID) {
			err = _error.HttpErrString("invalid user id", response.CM04)
		}
		helper.ErrorResponseEncode(w, err)
		return
	}

	resp := &schema.UserResponse{
		ID:              user.ID,
		FullName:        user.FullName,
		Gender:          user.Gender,
		Image:           user.Image,
		Username:        user.Username,
		Email:           user.Email,
		EmailFormat:     user.EmailFormat,
		PhoneNumber:     user.PhoneNumber,
		EmailVerifiedAt: user.EmailVerifiedAt,
	}
	helper.SuccessResponseEncode(w, resp, "data user")
}
