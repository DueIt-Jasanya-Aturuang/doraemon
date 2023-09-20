package unit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/mocks"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/api/rest"
)

func TestSecurityHandlerValidateAccess(t *testing.T) {
	securityUsecase := &mocks.FakeSecurityUsecase{}
	appUsecase := &mocks.FakeAppUsecase{}

	securityHandler := rest.NewSecurityHandlerImpl(securityUsecase, appUsecase)

	t.Run("SUCCESS", func(t *testing.T) {
		httpreq, err := http.NewRequest("POST", "/validate", nil)
		httpreq.Header.Set("App-ID", "appid")
		httpreq.Header.Set("User-ID", "userid")
		httpreq.Header.Set("Authorization", "token")
		assert.NoError(t, err)
		httpresp := httptest.NewRecorder()

		appUsecase.CheckAppByID(context.TODO(), &dto.AppReq{AppID: "appid"})
		appUsecase.CheckAppByIDReturns(nil)

		securityUsecase.JwtValidateAT(context.TODO(), &dto.JwtTokenReq{}, httpreq.URL.Path)
		securityUsecase.JwtValidateATReturns(false, nil)

		securityHandler.ValidateAccess(httpresp, httpreq)

		assert.Equal(t, 200, httpresp.Code)
		assert.Equal(t, "", httpresp.Header().Get("Authorization"))
	})

	t.Run("SUCCESS_with-generate-rt-at", func(t *testing.T) {
		httpreq, err := http.NewRequest("POST", "/validate", nil)
		httpreq.Header.Set("App-ID", "appid")
		httpreq.Header.Set("User-ID", "userid")
		httpreq.Header.Set("Authorization", "token")
		assert.NoError(t, err)
		httpresp := httptest.NewRecorder()

		appUsecase.CheckAppByID(context.TODO(), &dto.AppReq{AppID: "appid"})
		appUsecase.CheckAppByIDReturns(nil)

		securityUsecase.JwtValidateAT(context.TODO(), &dto.JwtTokenReq{}, httpreq.URL.Path)
		securityUsecase.JwtValidateATReturns(true, nil)

		securityUsecase.JwtGenerateRTAT(context.TODO(), &dto.JwtTokenReq{})
		securityUsecase.JwtGenerateRTATReturns(&dto.JwtTokenResp{Token: "asal"}, nil)

		securityHandler.ValidateAccess(httpresp, httpreq)

		assert.Equal(t, 200, httpresp.Code)
		assert.Equal(t, "asal", httpresp.Header().Get("Authorization"))
	})
}

func TestSecurityHandlerLogout(t *testing.T) {
	securityUsecase := &mocks.FakeSecurityUsecase{}
	appUsecase := &mocks.FakeAppUsecase{}

	securityHandler := rest.NewSecurityHandlerImpl(securityUsecase, appUsecase)

	t.Run("SUCCESS", func(t *testing.T) {
		httpreq, err := http.NewRequest("POST", "/validate", nil)
		httpreq.Header.Set("User-ID", "userid")
		httpreq.Header.Set("Authorization", "token")
		assert.NoError(t, err)
		httpresp := httptest.NewRecorder()

		securityUsecase.Logout(context.TODO(), &dto.LogoutReq{})
		securityUsecase.LogoutReturns(nil)

		securityHandler.Logout(httpresp, httpreq)

		assert.Equal(t, 200, httpresp.Code)
		t.Log(httpresp.Body)
	})
}
