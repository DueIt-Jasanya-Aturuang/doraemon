package unit

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/delivery/restapi"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/mocks"
)

func TestUserHandlerResetPassword(t *testing.T) {
	userUsecase := &mocks.FakeUserUsecase{}
	appUsecase := &mocks.FakeAppUsecase{}
	otpUsecase := &mocks.FakeOTPUsecase{}

	userHandler := restapi.NewUserHandlerImpl(userUsecase, appUsecase, otpUsecase)

	req := &dto.ResetPasswordReq{
		OldPassword: "",
		Password:    "",
		RePassword:  "",
		UserID:      "",
	}
	t.Run("SUCCESS", func(t *testing.T) {
		reqByte, err := json.Marshal(req)
		assert.NoError(t, err)

		httpreq, err := http.NewRequest("POST", "/reset-password", bytes.NewReader(reqByte))
		httpresp := httptest.NewRecorder()

		userHandler.ResetPassword(httpresp, httpreq)
	})
}

func TestUserHandlerForgottenPassword(t *testing.T) {
	userUsecase := &mocks.FakeUserUsecase{}
	appUsecase := &mocks.FakeAppUsecase{}
	otpUsecase := &mocks.FakeOTPUsecase{}

	userHandler := restapi.NewUserHandlerImpl(userUsecase, appUsecase, otpUsecase)

	req := &dto.OTPValidationReq{
		Email: "",
		OTP:   "",
	}
	t.Run("SUCCESS", func(t *testing.T) {
		reqByte, err := json.Marshal(req)
		assert.NoError(t, err)

		httpreq, err := http.NewRequest("POST", "/forgot-password", bytes.NewReader(reqByte))
		httpreq.Header.Set("Type", "forgot-password")
		httpresp := httptest.NewRecorder()

		userHandler.ForgottenPassword(httpresp, httpreq)
	})
}

func TestUserHandlerResetForgottenPassword(t *testing.T) {
	userUsecase := &mocks.FakeUserUsecase{}
	appUsecase := &mocks.FakeAppUsecase{}
	otpUsecase := &mocks.FakeOTPUsecase{}

	userHandler := restapi.NewUserHandlerImpl(userUsecase, appUsecase, otpUsecase)

	req := &dto.ResetForgottenPasswordReq{
		Password:   "",
		RePassword: "",
	}
	t.Run("SUCCESS", func(t *testing.T) {
		reqByte, err := json.Marshal(req)
		assert.NoError(t, err)

		httpreq, err := http.NewRequest("POST", "/reset-forgot-password?email=asd&token=asd", bytes.NewReader(reqByte))
		httpreq.Header.Set("Type", "forgot-password")
		httpresp := httptest.NewRecorder()

		userHandler.ResetForgottenPassword(httpresp, httpreq)
	})
}

func TestUserHandlerActivasiAccount(t *testing.T) {
	userUsecase := &mocks.FakeUserUsecase{}
	appUsecase := &mocks.FakeAppUsecase{}
	otpUsecase := &mocks.FakeOTPUsecase{}

	userHandler := restapi.NewUserHandlerImpl(userUsecase, appUsecase, otpUsecase)

	req := &dto.OTPValidationReq{
		Email: "",
		OTP:   "",
	}
	t.Run("SUCCESS", func(t *testing.T) {
		reqByte, err := json.Marshal(req)
		assert.NoError(t, err)

		httpreq, err := http.NewRequest("POST", "/activasi-account", bytes.NewReader(reqByte))
		httpreq.Header.Set("Type", "activasi-account")
		httpresp := httptest.NewRecorder()

		userHandler.ActivasiAccount(httpresp, httpreq)
	})
}
