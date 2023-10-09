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

	_error2 "github.com/DueIt-Jasanya-Aturuang/doraemon/util/error"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/old/rapi"
)

func TestAuthHandlerLogin(t *testing.T) {
	authUsecase := &mocks.FakeAuthUsecase{}
	securityUsecase := &mocks.FakeSecurityUsecase{}
	appUsecase := &mocks.FakeAppUsecase{}

	authHandler := rapi.NewAuthHandlerImpl(authUsecase, securityUsecase, appUsecase)

	req := dto.LoginReq{
		EmailOrUsername: "ibanrama29@gmail.com",
		Password:        "123456",
		RememberMe:      false,
		Oauth2:          false,
	}

	t.Run("SUCCESS_200", func(t *testing.T) {
		reqByte, err := json.Marshal(req)
		assert.NoError(t, err)

		appUsecase.CheckAppByID(context.TODO(), &dto.AppReq{AppID: "appid"})
		appUsecase.CheckAppByIDReturns(nil)

		authUsecase.Login(context.TODO(), &req)
		authUsecase.LoginReturns(&dto.UserResp{}, &dto.ProfileResp{}, nil)

		securityUsecase.JwtRegistredRTAT(context.TODO(), &dto.JwtRegisteredTokenReq{
			AppId:      "appid",
			UserId:     "userID_1",
			RememberMe: false,
		})
		securityUsecase.JwtRegistredRTATReturns(&dto.JwtTokenResp{}, nil)

		httpReq, err := http.NewRequest("POST", "/register", bytes.NewReader(reqByte))
		assert.NoError(t, err)
		httpReq.Header.Set("App-ID", "appid")

		responWriter := httptest.NewRecorder()
		authHandler.Login(responWriter, httpReq)

		assert.Equal(t, 200, responWriter.Code)
		t.Log(responWriter.Body)
	})

	t.Run("ERROR_403-appid-invalid", func(t *testing.T) {
		reqByte, err := json.Marshal(req)
		assert.NoError(t, err)

		httpReq, err := http.NewRequest("POST", "/register", bytes.NewReader(reqByte))
		assert.NoError(t, err)
		httpReq.Header.Set("App-ID", "appid")

		appUsecase.CheckAppByID(context.TODO(), &dto.AppReq{AppID: ""})
		appUsecase.CheckAppByIDReturns(_error2.ErrStringDefault(http.StatusForbidden))

		responWriter := httptest.NewRecorder()
		authHandler.Login(responWriter, httpReq)

		assert.Equal(t, 403, responWriter.Code)
		t.Log(responWriter.Body)
	})

	t.Run("ERROR_403-appid-nil", func(t *testing.T) {
		reqByte, err := json.Marshal(req)
		assert.NoError(t, err)

		httpReq, err := http.NewRequest("POST", "/register", bytes.NewReader(reqByte))
		assert.NoError(t, err)

		responWriter := httptest.NewRecorder()
		authHandler.Login(responWriter, httpReq)

		assert.Equal(t, 403, responWriter.Code)
		t.Log(responWriter.Body)
	})

	t.Run("ERROR_400-password-not-match", func(t *testing.T) {
		reqByte, err := json.Marshal(req)
		assert.NoError(t, err)

		httpReq, err := http.NewRequest("POST", "/register", bytes.NewReader(reqByte))
		assert.NoError(t, err)
		httpReq.Header.Set("App-ID", "appid")

		appUsecase.CheckAppByID(context.TODO(), &dto.AppReq{AppID: "appid"})
		appUsecase.CheckAppByIDReturns(nil)

		authUsecase.Login(context.TODO(), &req)
		authUsecase.LoginReturns(nil, nil, _error2.BadLogin())

		responWriter := httptest.NewRecorder()
		authHandler.Login(responWriter, httpReq)

		assert.Equal(t, 400, responWriter.Code)
		t.Log(responWriter.Body)
	})

	t.Run("ERROR_502-old-bad-gateway", func(t *testing.T) {
		reqByte, err := json.Marshal(req)
		assert.NoError(t, err)

		httpReq, err := http.NewRequest("POST", "/register", bytes.NewReader(reqByte))
		assert.NoError(t, err)
		httpReq.Header.Set("App-ID", "appid")

		appUsecase.CheckAppByID(context.TODO(), &dto.AppReq{AppID: "appid"})
		appUsecase.CheckAppByIDReturns(nil)

		authUsecase.Login(context.TODO(), &req)
		authUsecase.LoginReturns(nil, nil, _error2.ErrStringDefault(http.StatusBadGateway))

		responWriter := httptest.NewRecorder()
		authHandler.Login(responWriter, httpReq)

		assert.Equal(t, 502, responWriter.Code)
		t.Log(responWriter.Body)
	})

	t.Run("ERROR_400", func(t *testing.T) {
		req.Password = "123"
		req.EmailOrUsername = "123"
		reqByte, err := json.Marshal(req)
		assert.NoError(t, err)

		appUsecase.CheckAppByID(context.TODO(), &dto.AppReq{AppID: "appid"})
		appUsecase.CheckAppByIDReturns(nil)

		authUsecase.Login(context.TODO(), &req)
		authUsecase.LoginReturns(&dto.UserResp{}, &dto.ProfileResp{}, nil)

		securityUsecase.JwtRegistredRTAT(context.TODO(), &dto.JwtRegisteredTokenReq{
			AppId:      "appid",
			UserId:     "userID_1",
			RememberMe: false,
		})
		securityUsecase.JwtRegistredRTATReturns(&dto.JwtTokenResp{}, nil)

		httpReq, err := http.NewRequest("POST", "/register", bytes.NewReader(reqByte))
		assert.NoError(t, err)
		httpReq.Header.Set("App-ID", "appid")

		responWriter := httptest.NewRecorder()
		authHandler.Login(responWriter, httpReq)

		assert.Equal(t, 400, responWriter.Code)
		t.Log(responWriter.Body)
	})
}

func TestAuthHandlerRegister(t *testing.T) {
	authUsecase := &mocks.FakeAuthUsecase{}
	securityUsecase := &mocks.FakeSecurityUsecase{}
	appUsecase := &mocks.FakeAppUsecase{}

	authHandler := rapi.NewAuthHandlerImpl(authUsecase, securityUsecase, appUsecase)

	req := &dto.RegisterReq{
		FullName:   "ibanraa",
		Username:   "rama",
		Email:      "ibanrama29@gmail.com",
		Password:   "123456",
		RePassword: "123456",
	}

	t.Run("SUCCESS", func(t *testing.T) {
		reqByte, err := json.Marshal(req)
		assert.NoError(t, err)

		appUsecase.CheckAppByID(context.TODO(), &dto.AppReq{AppID: "appid"})
		appUsecase.CheckAppByIDReturns(nil)

		authUsecase.Register(context.TODO(), req)
		authUsecase.RegisterReturns(&dto.UserResp{}, nil)

		authUsecase.Login(context.TODO(), &dto.LoginReq{
			EmailOrUsername: req.Username,
			Password:        req.Password,
			RememberMe:      false,
			Oauth2:          false,
		})
		authUsecase.LoginReturns(&dto.UserResp{}, &dto.ProfileResp{}, nil)

		securityUsecase.JwtRegistredRTAT(context.TODO(), &dto.JwtRegisteredTokenReq{
			AppId:      "appid",
			UserId:     "userID",
			RememberMe: false,
		})
		securityUsecase.JwtRegistredRTATReturns(&dto.JwtTokenResp{}, nil)

		httpReq, err := http.NewRequest("POST", "/register", bytes.NewReader(reqByte))
		assert.NoError(t, err)
		httpReq.Header.Set("App-ID", "appid")

		responseWriter := httptest.NewRecorder()

		authHandler.Register(responseWriter, httpReq)

		assert.Equal(t, 200, responseWriter.Code)
		t.Log(responseWriter.Body)
	})

	t.Run("ERROR_appid-invalid", func(t *testing.T) {
		reqByte, err := json.Marshal(req)
		assert.NoError(t, err)

		appUsecase.CheckAppByID(context.TODO(), &dto.AppReq{AppID: "nil"})
		appUsecase.CheckAppByIDReturns(_error2.ErrStringDefault(http.StatusForbidden))

		httpReq, err := http.NewRequest("POST", "/register", bytes.NewReader(reqByte))
		assert.NoError(t, err)
		httpReq.Header.Set("App-ID", "appid")

		responseWriter := httptest.NewRecorder()

		authHandler.Register(responseWriter, httpReq)

		assert.Equal(t, 403, responseWriter.Code)
		t.Log(responseWriter.Body)
	})

	t.Run("ERROR_appid-nil", func(t *testing.T) {
		reqByte, err := json.Marshal(req)
		assert.NoError(t, err)

		httpReq, err := http.NewRequest("POST", "/register", bytes.NewReader(reqByte))
		assert.NoError(t, err)

		responseWriter := httptest.NewRecorder()

		authHandler.Register(responseWriter, httpReq)

		assert.Equal(t, 403, responseWriter.Code)
		t.Log(responseWriter.Body)
	})

	t.Run("ERROR_email-exist", func(t *testing.T) {
		reqByte, err := json.Marshal(req)
		assert.NoError(t, err)

		appUsecase.CheckAppByID(context.TODO(), &dto.AppReq{AppID: "appid"})
		appUsecase.CheckAppByIDReturns(nil)

		authUsecase.Register(context.TODO(), req)
		authUsecase.RegisterReturns(nil, _error2.BadExistField("email", "email has been registered"))

		httpReq, err := http.NewRequest("POST", "/register", bytes.NewReader(reqByte))
		assert.NoError(t, err)
		httpReq.Header.Set("App-ID", "appid")

		responseWriter := httptest.NewRecorder()

		authHandler.Register(responseWriter, httpReq)

		assert.Equal(t, 400, responseWriter.Code)
		t.Log(responseWriter.Body)
	})

	t.Run("ERROR_username-exist", func(t *testing.T) {
		reqByte, err := json.Marshal(req)
		assert.NoError(t, err)

		appUsecase.CheckAppByID(context.TODO(), &dto.AppReq{AppID: "appid"})
		appUsecase.CheckAppByIDReturns(nil)

		authUsecase.Register(context.TODO(), req)
		authUsecase.RegisterReturns(nil, _error2.BadExistField("username", "username has been registered"))

		httpReq, err := http.NewRequest("POST", "/register", bytes.NewReader(reqByte))
		assert.NoError(t, err)
		httpReq.Header.Set("App-ID", "appid")

		responseWriter := httptest.NewRecorder()

		authHandler.Register(responseWriter, httpReq)

		assert.Equal(t, 400, responseWriter.Code)
		t.Log(responseWriter.Body)
	})

	t.Run("ERROR_badgateway-old-account", func(t *testing.T) {
		reqByte, err := json.Marshal(req)
		assert.NoError(t, err)

		appUsecase.CheckAppByID(context.TODO(), &dto.AppReq{AppID: "appid"})
		appUsecase.CheckAppByIDReturns(nil)

		authUsecase.Register(context.TODO(), req)
		authUsecase.RegisterReturns(nil, _error2.ErrStringDefault(http.StatusBadGateway))

		httpReq, err := http.NewRequest("POST", "/register", bytes.NewReader(reqByte))
		assert.NoError(t, err)
		httpReq.Header.Set("App-ID", "appid")

		responseWriter := httptest.NewRecorder()

		authHandler.Register(responseWriter, httpReq)

		assert.Equal(t, 502, responseWriter.Code)
		t.Log(responseWriter.Body)
	})

	t.Run("ERROR_400-password-not-match", func(t *testing.T) {
		reqByte, err := json.Marshal(req)
		assert.NoError(t, err)

		httpReq, err := http.NewRequest("POST", "/register", bytes.NewReader(reqByte))
		assert.NoError(t, err)
		httpReq.Header.Set("App-ID", "appid")

		appUsecase.CheckAppByID(context.TODO(), &dto.AppReq{AppID: "appid"})
		appUsecase.CheckAppByIDReturns(nil)

		authUsecase.Register(context.TODO(), req)
		authUsecase.RegisterReturns(&dto.UserResp{}, nil)

		authUsecase.Login(context.TODO(), &dto.LoginReq{
			EmailOrUsername: req.Username,
			Password:        req.Password,
			RememberMe:      false,
			Oauth2:          false,
		})
		authUsecase.LoginReturns(nil, nil, _error2.BadLogin())

		responWriter := httptest.NewRecorder()
		authHandler.Register(responWriter, httpReq)

		assert.Equal(t, 400, responWriter.Code)
		t.Log(responWriter.Body)
	})

	t.Run("ERROR_api-badgateway-login", func(t *testing.T) {
		reqByte, err := json.Marshal(req)
		assert.NoError(t, err)

		httpReq, err := http.NewRequest("POST", "/register", bytes.NewReader(reqByte))
		assert.NoError(t, err)
		httpReq.Header.Set("App-ID", "appid")

		appUsecase.CheckAppByID(context.TODO(), &dto.AppReq{AppID: "appid"})
		appUsecase.CheckAppByIDReturns(nil)

		authUsecase.Register(context.TODO(), req)
		authUsecase.RegisterReturns(&dto.UserResp{}, nil)

		authUsecase.Login(context.TODO(), &dto.LoginReq{
			EmailOrUsername: req.Username,
			Password:        req.Password,
			RememberMe:      false,
			Oauth2:          false,
		})
		authUsecase.LoginReturns(nil, nil, _error2.ErrStringDefault(http.StatusBadGateway))

		responWriter := httptest.NewRecorder()
		authHandler.Register(responWriter, httpReq)

		assert.Equal(t, 502, responWriter.Code)
		t.Log(responWriter.Body)
	})

	t.Run("ERROR_bad-request-input", func(t *testing.T) {
		req.Password = "123"
		req.RePassword = "1234"
		req.Email = "invalid"
		req.Username = "12"
		req.FullName = "12"

		reqByte, err := json.Marshal(req)
		assert.NoError(t, err)

		httpReq, err := http.NewRequest("POST", "/register", bytes.NewReader(reqByte))
		assert.NoError(t, err)
		httpReq.Header.Set("App-ID", "appid")

		appUsecase.CheckAppByID(context.TODO(), &dto.AppReq{AppID: "appid"})
		appUsecase.CheckAppByIDReturns(nil)

		responWriter := httptest.NewRecorder()
		authHandler.Register(responWriter, httpReq)

		assert.Equal(t, 400, responWriter.Code)
		t.Log(responWriter.Body)
	})
}
