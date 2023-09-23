package rest

import (
	"errors"
	"net/http"

	"github.com/jasanya-tech/jasanya-response-backend-golang/_error"
	"github.com/jasanya-tech/jasanya-response-backend-golang/response"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/api/rest/helper"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/api/validation"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/pkg/_usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

type AuthHandlerImpl struct {
	authUsecase domain.AuthUsecase
	otpUsecase  domain.OTPUsecase
}

func NewAuthHandlerImpl(
	authUsecase domain.AuthUsecase,
	otpUsecase domain.OTPUsecase,
) *AuthHandlerImpl {
	return &AuthHandlerImpl{
		authUsecase: authUsecase,
		otpUsecase:  otpUsecase,
	}
}

func (h *AuthHandlerImpl) Register(w http.ResponseWriter, r *http.Request) {
	// with app middleware

	// decod request ke dalam dto
	req := new(domain.RequestRegister)
	err := helper.DecodeJson(r, &req)
	if err != nil {
		helper.ErrorResponseEncode(w, err)
		return
	}

	req.Role = 1
	req.AppID = r.Header.Get(util.AppIDHeader)
	req.EmailVerifiedAt = false

	err = validation.RegisterValidation(req)
	if err != nil {
		helper.ErrorResponseEncode(w, err)
		return
	}

	// register user
	resp, err := h.authUsecase.Register(r.Context(), req)
	if err != nil {
		if errors.Is(err, _usecase.EmailIsExist) {
			err = _error.HttpErrMapOfSlices(map[string][]string{
				"email": {
					"email sudah terdaftar",
				},
			}, response.CM06)
		}
		if errors.Is(err, _usecase.UsernameIsExist) {
			err = _error.HttpErrMapOfSlices(map[string][]string{
				"username": {
					"username sudah tersedia",
				},
			}, response.CM06)
		}
		helper.ErrorResponseEncode(w, err)
		return
	}

	// mencoba untuk push otp langsung, ga peduli error atau tidak
	_ = h.otpUsecase.Generate(r.Context(), &domain.RequestGenerateOTP{
		Email:  req.Email,
		Type:   util.ActivasiAccount,
		UserID: resp.ID,
	})

	helper.SuccessResponseEncode(w, req, "register successfully")
}

func (h *AuthHandlerImpl) Login(w http.ResponseWriter, r *http.Request) {

	req := new(domain.RequestLogin)
	err := helper.DecodeJson(r, &req)
	if err != nil {
		helper.ErrorResponseEncode(w, err)
		return
	}
	req.AppID = r.Header.Get(util.AppIDHeader)

	err = validation.LoginValidation(req)
	if err != nil {
		helper.ErrorResponseEncode(w, err)
		return
	}

	resp, err := h.authUsecase.Login(r.Context(), req)
	if err != nil {
		if errors.Is(err, _usecase.InvalidEmailOrUsernameOrPassword) {
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

	helper.SuccessResponseEncode(w, resp, "login successfully")
}
