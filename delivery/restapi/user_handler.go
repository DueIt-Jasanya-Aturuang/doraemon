package restapi

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/delivery/restapi/mapper"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/delivery/validation"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/usecase"
	_error "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/error"
)

type UserHandlerImpl struct {
	userUsecase usecase.UserUsecase
	appUsecase  usecase.AppUsecase
	otpUsecase  usecase.OTPUsecase
}

func NewUserHandlerImpl(
	userUsecase usecase.UserUsecase,
	appUsecase usecase.AppUsecase,
	otpUsecase usecase.OTPUsecase,
) *UserHandlerImpl {
	return &UserHandlerImpl{
		userUsecase: userUsecase,
		appUsecase:  appUsecase,
		otpUsecase:  otpUsecase,
	}
}

func (h *UserHandlerImpl) ResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var reqResetPassword dto.ResetPasswordReq
	err := mapper.DecodeJson(r, &reqResetPassword)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	userID := r.Header.Get("User-ID")
	reqResetPassword.UserID = userID

	err = validation.ResetPasswordValidation(&reqResetPassword)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	err = h.userUsecase.ResetPassword(ctx, &reqResetPassword)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	resp := mapper.ResponseSuccess{
		Message: "password anda telah berhasil dirubah",
	}

	mapper.NewSuccessResp(w, r, resp, 200)

}

func (h *UserHandlerImpl) ForgottenPassword(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	appID := r.Header.Get("App-ID")
	if appID == "" {
		log.Warn().Msgf("app id header tidak tersedia")
		mapper.NewErrorResp(w, r, _error.ErrStringDefault(http.StatusForbidden))
		return
	}

	err := h.appUsecase.CheckAppByID(ctx, &dto.AppReq{
		AppID: appID,
	})
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	var reqOtpValidation dto.OTPValidationReq
	err = mapper.DecodeJson(r, &reqOtpValidation)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}
	reqOtpValidation.Type = "forgot-password"

	err = validation.OTPValidation(&reqOtpValidation)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	err = h.otpUsecase.OTPValidation(ctx, &reqOtpValidation)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	var reqForgottenPassword dto.ForgottenPasswordReq
	reqForgottenPassword.Email = reqOtpValidation.Email

	url, err := h.userUsecase.ForgottenPassword(ctx, &reqForgottenPassword)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	resp := mapper.ResponseSuccess{
		Data: map[string]string{
			"url_forgot_password": url,
		},
	}

	mapper.NewSuccessResp(w, r, resp, 200)
}

func (h *UserHandlerImpl) ResetForgottenPassword(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	appID := r.Header.Get("App-ID")
	if appID == "" {
		log.Warn().Msgf("app id header tidak tersedia")
		mapper.NewErrorResp(w, r, _error.ErrStringDefault(http.StatusForbidden))
		return
	}

	err := h.appUsecase.CheckAppByID(ctx, &dto.AppReq{
		AppID: appID,
	})
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	var reqResetForgottenPassword dto.ResetForgottenPasswordReq
	err = mapper.DecodeJson(r, &reqResetForgottenPassword)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}
	email := r.URL.Query().Get("email")
	token := r.URL.Query().Get("token")
	reqResetForgottenPassword.Email = email
	reqResetForgottenPassword.Token = token

	err = validation.ResetForgottenPasswordValidation(&reqResetForgottenPassword)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	err = h.userUsecase.ResetForgottenPassword(ctx, &reqResetForgottenPassword)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	resp := mapper.ResponseSuccess{
		Message: "password anda telah berhasil diubah",
	}

	mapper.NewSuccessResp(w, r, resp, 200)
}

func (h *UserHandlerImpl) ActivasiAccount(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var reqValidationOtp dto.OTPValidationReq
	err := mapper.DecodeJson(r, &reqValidationOtp)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}
	reqValidationOtp.Type = "activasi-account"

	err = validation.OTPValidation(&reqValidationOtp)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	err = h.otpUsecase.OTPValidation(ctx, &reqValidationOtp)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	activasi, err := h.userUsecase.ActivasiAccount(ctx, reqValidationOtp.Email)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	resp := mapper.ResponseSuccess{
		Data: activasi,
	}

	mapper.NewSuccessResp(w, r, resp, 200)
}
