package unit

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/mocks"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/api/rest"
)

func TestOTPHandlerOTPGenerate(t *testing.T) {
	otpUsecase := &mocks.FakeOTPUsecase{}

	otpHandler := rest.NewOTPHandlerImpl(otpUsecase)

	t.Run("SUCCESS_activasi-account", func(t *testing.T) {
		req := &dto.OTPGenerateReq{
			Email: "ibanrama29@gmail.com",
		}
		bytesReq, err := json.Marshal(req)
		assert.NoError(t, err)

		otpUsecase.OTPGenerate(context.TODO(), req)
		otpUsecase.OTPGenerateReturns(nil)

		httpreq, err := http.NewRequest("POST", "/generateotp", bytes.NewReader(bytesReq))
		httpreq.Header.Set("User-ID", "123")
		httpreq.Header.Set("Type", "activasi-account")
		httpResp := httptest.NewRecorder()

		otpHandler.GenerateOTP(httpResp, httpreq)

		assert.Equal(t, 200, httpResp.Code)
		t.Log(httpResp.Body)
	})

	t.Run("SUCCESS_forgot-password", func(t *testing.T) {
		req := &dto.OTPGenerateReq{
			Email: "ibanrama29@gmail.com",
		}
		bytesReq, err := json.Marshal(req)
		assert.NoError(t, err)

		otpUsecase.OTPGenerate(context.TODO(), req)
		otpUsecase.OTPGenerateReturns(nil)

		httpreq, err := http.NewRequest("POST", "/generateotp", bytes.NewReader(bytesReq))
		httpreq.Header.Set("Type", "forgot-password")
		httpResp := httptest.NewRecorder()

		otpHandler.GenerateOTP(httpResp, httpreq)

		assert.Equal(t, 200, httpResp.Code)
		t.Log(httpResp.Body)
	})

	t.Run("ERROR_activasi-account-activasi-true", func(t *testing.T) {
		req := &dto.OTPGenerateReq{
			Email: "ibanrama29@gmail.com",
		}
		bytesReq, err := json.Marshal(req)
		assert.NoError(t, err)

		otpUsecase.OTPGenerate(context.TODO(), req)
		otpUsecase.OTPGenerateReturns(_error.Err400(map[string][]string{
			"email": {
				"permintaan anda tidak dapat di proses, email anda sudah di aktivasi silahkan login",
			},
		}))

		httpreq, err := http.NewRequest("POST", "/generateotp", bytes.NewReader(bytesReq))
		httpreq.Header.Set("User-ID", "123")
		httpreq.Header.Set("Type", "activasi-account")
		httpResp := httptest.NewRecorder()

		otpHandler.GenerateOTP(httpResp, httpreq)

		assert.Equal(t, 400, httpResp.Code)
		t.Log(httpResp.Body)
	})

	t.Run("ERROR_activasi-account-invalid-userid", func(t *testing.T) {
		req := &dto.OTPGenerateReq{
			Email: "ibanrama29@gmail.com",
		}
		bytesReq, err := json.Marshal(req)
		assert.NoError(t, err)

		otpUsecase.OTPGenerate(context.TODO(), req)
		otpUsecase.OTPGenerateReturns(nil)

		httpreq, err := http.NewRequest("POST", "/generateotp", bytes.NewReader(bytesReq))
		httpreq.Header.Set("Type", "activasi-account")
		httpResp := httptest.NewRecorder()

		otpHandler.GenerateOTP(httpResp, httpreq)

		assert.Equal(t, 403, httpResp.Code)
		t.Log(httpResp.Body)
	})

	t.Run("ERROR_activasi-account-email-notavailable", func(t *testing.T) {
		req := &dto.OTPGenerateReq{
			Email: "ibanrama29@gmail.com",
		}
		bytesReq, err := json.Marshal(req)
		assert.NoError(t, err)

		otpUsecase.OTPGenerate(context.TODO(), req)
		otpUsecase.OTPGenerateReturns(_error.ErrStringDefault(http.StatusNotFound))

		httpreq, err := http.NewRequest("POST", "/generateotp", bytes.NewReader(bytesReq))
		httpreq.Header.Set("User-ID", "123")
		httpreq.Header.Set("Type", "activasi-account")
		httpResp := httptest.NewRecorder()

		otpHandler.GenerateOTP(httpResp, httpreq)

		assert.Equal(t, 404, httpResp.Code)
		t.Log(httpResp.Body)
	})

	t.Run("ERROR_forgot-password-email-notavailable", func(t *testing.T) {
		req := &dto.OTPGenerateReq{
			Email: "ibanrama29@gmail.com",
		}
		bytesReq, err := json.Marshal(req)
		assert.NoError(t, err)

		otpUsecase.OTPGenerate(context.TODO(), req)
		otpUsecase.OTPGenerateReturns(_error.ErrStringDefault(http.StatusNotFound))

		httpreq, err := http.NewRequest("POST", "/generateotp", bytes.NewReader(bytesReq))
		httpreq.Header.Set("Type", "forgot-password")
		httpResp := httptest.NewRecorder()

		otpHandler.GenerateOTP(httpResp, httpreq)

		assert.Equal(t, 404, httpResp.Code)
		t.Log(httpResp.Body)
	})

	t.Run("ERROR_forgot-password-invalid-type", func(t *testing.T) {
		req := &dto.OTPGenerateReq{
			Email: "ibanrama29@gmail.com",
		}
		bytesReq, err := json.Marshal(req)
		assert.NoError(t, err)

		httpreq, err := http.NewRequest("POST", "/generateotp", bytes.NewReader(bytesReq))
		httpreq.Header.Set("Type", "forgt-password")
		httpResp := httptest.NewRecorder()

		otpHandler.GenerateOTP(httpResp, httpreq)

		assert.Equal(t, 403, httpResp.Code)
		t.Log(httpResp.Body)
	})

	t.Run("ERROR_activasi-account-invalid-type", func(t *testing.T) {
		req := &dto.OTPGenerateReq{
			Email: "ibanrama29@gmail.com",
		}
		bytesReq, err := json.Marshal(req)
		assert.NoError(t, err)

		httpreq, err := http.NewRequest("POST", "/generateotp", bytes.NewReader(bytesReq))
		httpreq.Header.Set("Type", "activasi-accounts")
		httpResp := httptest.NewRecorder()

		otpHandler.GenerateOTP(httpResp, httpreq)

		assert.Equal(t, 403, httpResp.Code)
		t.Log(httpResp.Body)
	})
}
