package restapi

import (
	"context"
	"net/http"
	"time"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/delivery/restapi/mapper"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/delivery/validation"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/usecase"
	_error "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/error"
)

type AuthHandlerImpl struct {
	authUsecase     usecase.AuthUsecase
	securityUsecase usecase.SecurityUsecase
	appUsecase      usecase.AppUsecase
}

func NewAuthHandlerImpl(
	authUsecase usecase.AuthUsecase,
	securityUsecase usecase.SecurityUsecase,
	appUsecase usecase.AppUsecase,
) *AuthHandlerImpl {
	return &AuthHandlerImpl{
		authUsecase:     authUsecase,
		securityUsecase: securityUsecase,
		appUsecase:      appUsecase,
	}
}

func (h *AuthHandlerImpl) Register(w http.ResponseWriter, r *http.Request) {
	var reqRegister dto.RegisterReq

	err := mapper.DecodeJson(r, &reqRegister)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	reqRegister.Role = 1
	reqRegister.AppID = r.Header.Get("App-ID")
	reqRegister.EmailVerifiedAt = false

	err = validation.RegisterValidation(&reqRegister)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	err = h.appUsecase.CheckAppByID(ctx, &dto.AppReq{
		AppID: reqRegister.AppID,
	})
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	userResp, err := h.authUsecase.Register(ctx, &reqRegister)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	reqLogin := &dto.LoginReq{
		EmailOrUsername: userResp.Email,
		Password:        reqRegister.Password,
		RememberMe:      false,
		Oauth2:          false,
	}
	userResp, profileResp, err := h.authUsecase.Login(ctx, reqLogin)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	reqJwtRegistered := &dto.JwtRegisteredTokenReq{
		AppId:      reqRegister.AppID,
		UserId:     userResp.ID,
		RememberMe: reqLogin.RememberMe,
	}
	token, err := h.securityUsecase.JwtRegistredRTAT(ctx, reqJwtRegistered)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	respSuccess := mapper.ResponseSuccess{
		Data: map[string]any{
			"user":    userResp,
			"profile": profileResp,
			"token":   token,
		},
	}

	mapper.NewSuccessResp(w, r, respSuccess, 200)
}

func (h *AuthHandlerImpl) Login(w http.ResponseWriter, r *http.Request) {
	var reqLogin dto.LoginReq

	err := mapper.DecodeJson(r, &reqLogin)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	appID := r.Header.Get("App-ID")
	if appID == "" {
		mapper.NewErrorResp(w, r, _error.ErrStringDefault(http.StatusForbidden))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	err = h.appUsecase.CheckAppByID(ctx, &dto.AppReq{
		AppID: appID,
	})
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	userResp, profileResp, err := h.authUsecase.Login(ctx, &reqLogin)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	token, err := h.securityUsecase.JwtRegistredRTAT(ctx, &dto.JwtRegisteredTokenReq{
		AppId:      appID,
		UserId:     userResp.ID,
		RememberMe: reqLogin.RememberMe,
	})
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	resp := mapper.ResponseSuccess{
		Data: map[string]any{
			"user":    userResp,
			"profile": profileResp,
			"token":   token,
		},
	}

	mapper.NewSuccessResp(w, r, resp, 200)
}
