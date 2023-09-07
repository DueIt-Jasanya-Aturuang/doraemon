package unit

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/delivery/restapi"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/mocks"
	_error "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/error"
)

func TestOauth2HandlerLoginWithGoogle(t *testing.T) {
	oauth2Usecase := &mocks.FakeOauth2Usecase{}
	authUsecase := &mocks.FakeAuthUsecase{}
	securityUsecase := &mocks.FakeSecurityUsecase{}
	appUsecase := &mocks.FakeAppUsecase{}

	oauth2Handler := restapi.NewOauth2HandlerImpl(oauth2Usecase, authUsecase, securityUsecase, appUsecase)

	req := &dto.LoginGoogleReq{
		Token:  "this tokenss",
		Device: "web",
	}

	t.Run("SUCCESS_register", func(t *testing.T) {
		reqBytes, err := json.Marshal(req)
		assert.NoError(t, err)

		appUsecase.CheckAppByID(context.TODO(), &dto.AppReq{AppID: "appid"})
		appUsecase.CheckAppByIDReturns(nil)

		oauth2Usecase.GoogleClaimUser(context.TODO(), req)
		oauth2Usecase.GoogleClaimUserReturns(&dto.LoginGoogleResp{}, nil)

		authUsecase.Register(context.TODO(), &dto.RegisterReq{})
		authUsecase.RegisterReturns(&dto.UserResp{}, nil)

		authUsecase.Login(context.TODO(), &dto.LoginReq{})
		authUsecase.LoginReturns(&dto.UserResp{}, &dto.ProfileResp{}, nil)

		securityUsecase.JwtRegistredRTAT(context.TODO(), &dto.JwtRegisteredTokenReq{})
		securityUsecase.JwtGenerateRTATReturns(&dto.JwtTokenResp{}, nil)

		httpReq, err := http.NewRequest("POST", "/login-google", bytes.NewReader(reqBytes))
		httpReq.Header.Set("App-ID", "appid")
		responWriter := httptest.NewRecorder()

		oauth2Handler.LoginWithGoogle(responWriter, httpReq)

		assert.Equal(t, 200, responWriter.Code)
		t.Log(responWriter.Body)
	})

	t.Run("SUCCESS_not-register", func(t *testing.T) {
		reqBytes, err := json.Marshal(req)
		assert.NoError(t, err)

		appUsecase.CheckAppByID(context.TODO(), &dto.AppReq{AppID: "appid"})
		appUsecase.CheckAppByIDReturns(nil)

		oauth2Usecase.GoogleClaimUser(context.TODO(), req)
		oauth2Usecase.GoogleClaimUserReturns(&dto.LoginGoogleResp{ExistsUser: true}, nil)

		authUsecase.Register(context.TODO(), &dto.RegisterReq{})
		authUsecase.RegisterReturns(&dto.UserResp{}, nil)

		authUsecase.Login(context.TODO(), &dto.LoginReq{})
		authUsecase.LoginReturns(&dto.UserResp{}, &dto.ProfileResp{}, nil)

		securityUsecase.JwtRegistredRTAT(context.TODO(), &dto.JwtRegisteredTokenReq{})
		securityUsecase.JwtGenerateRTATReturns(&dto.JwtTokenResp{}, nil)

		httpReq, err := http.NewRequest("POST", "/login-google", bytes.NewReader(reqBytes))
		httpReq.Header.Set("App-ID", "appid")
		responWriter := httptest.NewRecorder()

		oauth2Handler.LoginWithGoogle(responWriter, httpReq)

		assert.Equal(t, 200, responWriter.Code)
		t.Log(responWriter.Body)
	})

	t.Run("ERROR_appid-nil", func(t *testing.T) {
		reqBytes, err := json.Marshal(req)
		assert.NoError(t, err)

		httpReq, err := http.NewRequest("POST", "/login-google", bytes.NewReader(reqBytes))
		responWriter := httptest.NewRecorder()

		oauth2Handler.LoginWithGoogle(responWriter, httpReq)

		assert.Equal(t, 403, responWriter.Code)
		t.Log(responWriter.Body)
	})

	t.Run("ERROR_appid-invalid", func(t *testing.T) {
		reqBytes, err := json.Marshal(req)
		assert.NoError(t, err)

		appUsecase.CheckAppByID(context.TODO(), &dto.AppReq{AppID: "appid"})
		appUsecase.CheckAppByIDReturns(_error.ErrStringDefault(http.StatusForbidden))

		httpReq, err := http.NewRequest("POST", "/login-google", bytes.NewReader(reqBytes))
		httpReq.Header.Set("App-ID", "appid")
		responWriter := httptest.NewRecorder()

		oauth2Handler.LoginWithGoogle(responWriter, httpReq)

		assert.Equal(t, 403, responWriter.Code)
		t.Log(responWriter.Body)
	})

	t.Run("ERROR_claim-user-invalid-token", func(t *testing.T) {
		reqBytes, err := json.Marshal(req)
		assert.NoError(t, err)

		appUsecase.CheckAppByID(context.TODO(), &dto.AppReq{AppID: "appid"})
		appUsecase.CheckAppByIDReturns(nil)

		oauth2Usecase.GoogleClaimUser(context.TODO(), req)
		oauth2Usecase.GoogleClaimUserReturns(nil, _error.ErrStringDefault(http.StatusForbidden))

		httpReq, err := http.NewRequest("POST", "/login-google", bytes.NewReader(reqBytes))
		httpReq.Header.Set("App-ID", "appid")
		responWriter := httptest.NewRecorder()

		oauth2Handler.LoginWithGoogle(responWriter, httpReq)

		assert.Equal(t, 403, responWriter.Code)
		t.Log(responWriter.Body)
	})

	t.Run("ERROR_when-register", func(t *testing.T) {
		reqBytes, err := json.Marshal(req)
		assert.NoError(t, err)

		appUsecase.CheckAppByID(context.TODO(), &dto.AppReq{AppID: "appid"})
		appUsecase.CheckAppByIDReturns(nil)

		oauth2Usecase.GoogleClaimUser(context.TODO(), req)
		oauth2Usecase.GoogleClaimUserReturns(&dto.LoginGoogleResp{ExistsUser: false}, nil)

		authUsecase.Register(context.TODO(), &dto.RegisterReq{})
		authUsecase.RegisterReturns(nil, _error.ErrStringDefault(http.StatusInternalServerError))

		httpReq, err := http.NewRequest("POST", "/login-google", bytes.NewReader(reqBytes))
		httpReq.Header.Set("App-ID", "appid")
		responWriter := httptest.NewRecorder()

		oauth2Handler.LoginWithGoogle(responWriter, httpReq)

		assert.Equal(t, 500, responWriter.Code)
		t.Log(responWriter.Body)
	})

	t.Run("ERROR_when-login-password-notmatch", func(t *testing.T) {
		reqBytes, err := json.Marshal(req)
		assert.NoError(t, err)

		appUsecase.CheckAppByID(context.TODO(), &dto.AppReq{AppID: "appid"})
		appUsecase.CheckAppByIDReturns(nil)

		oauth2Usecase.GoogleClaimUser(context.TODO(), req)
		oauth2Usecase.GoogleClaimUserReturns(&dto.LoginGoogleResp{ExistsUser: true}, nil)

		authUsecase.Register(context.TODO(), &dto.RegisterReq{})
		authUsecase.RegisterReturns(&dto.UserResp{}, nil)

		authUsecase.Login(context.TODO(), &dto.LoginReq{})
		authUsecase.LoginReturns(nil, nil, _error.BadLogin())

		httpReq, err := http.NewRequest("POST", "/login-google", bytes.NewReader(reqBytes))
		httpReq.Header.Set("App-ID", "appid")
		responWriter := httptest.NewRecorder()

		oauth2Handler.LoginWithGoogle(responWriter, httpReq)

		assert.Equal(t, 400, responWriter.Code)
		t.Log(responWriter.Body)
	})

	t.Run("ERROR_when-login-badgateway", func(t *testing.T) {
		reqBytes, err := json.Marshal(req)
		assert.NoError(t, err)

		appUsecase.CheckAppByID(context.TODO(), &dto.AppReq{AppID: "appid"})
		appUsecase.CheckAppByIDReturns(nil)

		oauth2Usecase.GoogleClaimUser(context.TODO(), req)
		oauth2Usecase.GoogleClaimUserReturns(&dto.LoginGoogleResp{ExistsUser: true}, nil)

		authUsecase.Register(context.TODO(), &dto.RegisterReq{})
		authUsecase.RegisterReturns(&dto.UserResp{}, nil)

		authUsecase.Login(context.TODO(), &dto.LoginReq{})
		authUsecase.LoginReturns(nil, nil, _error.ErrStringDefault(http.StatusBadGateway))

		httpReq, err := http.NewRequest("POST", "/login-google", bytes.NewReader(reqBytes))
		httpReq.Header.Set("App-ID", "appid")
		responWriter := httptest.NewRecorder()

		oauth2Handler.LoginWithGoogle(responWriter, httpReq)

		assert.Equal(t, 502, responWriter.Code)
		t.Log(responWriter.Body)
	})

	t.Run("ERROR_badreq-invalid-input", func(t *testing.T) {
		reqNew := &dto.LoginGoogleReq{
			Token:  "asd",
			Device: "wseb",
		}

		reqBytes, err := json.Marshal(reqNew)
		assert.NoError(t, err)

		appUsecase.CheckAppByID(context.TODO(), &dto.AppReq{AppID: "appid"})
		appUsecase.CheckAppByIDReturns(nil)

		httpReq, err := http.NewRequest("POST", "/login-google", bytes.NewReader(reqBytes))
		httpReq.Header.Set("App-ID", "appid")
		responWriter := httptest.NewRecorder()

		oauth2Handler.LoginWithGoogle(responWriter, httpReq)

		assert.Equal(t, 403, responWriter.Code)
		t.Log(responWriter.Body)
	})
}
