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
	// set time out proccess
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// decode request body kedalam dto
	var reqResetPassword dto.ResetPasswordReq
	err := mapper.DecodeJson(r, &reqResetPassword)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	// get request userid di header
	userID := r.Header.Get("User-ID")
	reqResetPassword.UserID = userID

	// validasi request
	err = validation.ResetPasswordValidation(&reqResetPassword)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	// reset password process
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
	// set time out proccess
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// get header App-ID
	// jika gak ada akan return forbidden
	appID := r.Header.Get("App-ID")
	if appID == "" {
		log.Warn().Msgf("app id header tidak tersedia")
		mapper.NewErrorResp(w, r, _error.ErrStringDefault(http.StatusForbidden))
		return
	}

	// check appid, jika error akan return error
	// ini error sudah di set dari usecase, apakah error tersebut 500 atau yang lainnya
	err := h.appUsecase.CheckAppByID(ctx, &dto.AppReq{
		AppID: appID,
	})
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	// decod request ke dalam dto
	// dan set request otp type nya forgot-password
	var reqOtpValidation dto.OTPValidationReq
	err = mapper.DecodeJson(r, &reqOtpValidation)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}
	reqOtpValidation.Type = "forgot-password"

	// validasi request
	err = validation.OTPValidation(&reqOtpValidation)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	// process validasi otp
	err = h.otpUsecase.OTPValidation(ctx, &reqOtpValidation)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	// declare request forgotten password
	var reqForgottenPassword dto.ForgottenPasswordReq
	reqForgottenPassword.Email = reqOtpValidation.Email

	// process forgotten password
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
	// set time out proccess
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// get header App-ID
	// jika gak ada akan return forbidden
	appID := r.Header.Get("App-ID")
	if appID == "" {
		log.Warn().Msgf("app id header tidak tersedia")
		mapper.NewErrorResp(w, r, _error.ErrStringDefault(http.StatusForbidden))
		return
	}

	// check appid, jika error akan return error
	// ini error sudah di set dari usecase, apakah error tersebut 500 atau yang lainnya
	err := h.appUsecase.CheckAppByID(ctx, &dto.AppReq{
		AppID: appID,
	})
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	// decode request user
	var reqResetForgottenPassword dto.ResetForgottenPasswordReq
	err = mapper.DecodeJson(r, &reqResetForgottenPassword)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	// declare variable from url query param
	email := r.URL.Query().Get("email")
	token := r.URL.Query().Get("token")
	reqResetForgottenPassword.Email = email
	reqResetForgottenPassword.Token = token

	// validasi request
	err = validation.ResetForgottenPasswordValidation(&reqResetForgottenPassword)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	// process reset forgotten password
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
	// set time out proccess
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// decod request ke dalam dto
	// dan set type otp
	var reqValidationOtp dto.OTPValidationReq
	err := mapper.DecodeJson(r, &reqValidationOtp)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}
	reqValidationOtp.Type = "activasi-account"

	// validasi request otp
	err = validation.OTPValidation(&reqValidationOtp)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	// validasi otp
	err = h.otpUsecase.OTPValidation(ctx, &reqValidationOtp)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	// process activasi account
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
