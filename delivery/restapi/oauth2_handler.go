package restapi

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/delivery/restapi/mapper"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/delivery/validation"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/usecase"
	_error "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/error"
)

type Oauth2HandlerImpl struct {
	oauth2Usecase   usecase.Oauth2Usecase
	authUsecase     usecase.AuthUsecase
	securityUsecase usecase.SecurityUsecase
	appUsecase      usecase.AppUsecase
}

func NewOauth2HandlerImpl(
	oauth2Usecase usecase.Oauth2Usecase,
	authUsecase usecase.AuthUsecase,
	securityUsecase usecase.SecurityUsecase,
	appUsecase usecase.AppUsecase,
) *Oauth2HandlerImpl {
	return &Oauth2HandlerImpl{
		oauth2Usecase:   oauth2Usecase,
		authUsecase:     authUsecase,
		securityUsecase: securityUsecase,
		appUsecase:      appUsecase,
	}
}

func (h *Oauth2HandlerImpl) LoginWithGoogle(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	appID := r.Header.Get("App-ID")
	if appID == "" {
		log.Warn().Msg("tidak ada header appid")
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

	var reqLogin dto.LoginGoogleReq

	err = mapper.DecodeJson(r, &reqLogin)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	err = validation.Oauth2LoginValidation(&reqLogin)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	userGoogle, err := h.oauth2Usecase.GoogleClaimUser(ctx, &reqLogin)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	if !userGoogle.ExistsUser {
		reqRegister := &dto.RegisterReq{
			FullName:        userGoogle.Name,
			Username:        userGoogle.GivenName,
			Email:           userGoogle.Email,
			Password:        userGoogle.ID,
			RePassword:      userGoogle.ID,
			EmailVerifiedAt: true,
			AppID:           appID,
			Role:            1,
		}
		_, err = h.authUsecase.Register(ctx, reqRegister)
		if err != nil {
			var errHTTP *model.ErrResponseHTTP
			ok := errors.As(err, &errHTTP)
			if !ok {
				mapper.NewErrorResp(w, r, err)
				return
			}
			if errHTTP.Code == 500 || errHTTP.Code == 502 {
				mapper.NewErrorResp(w, r, err)
				return
			}
		}
	}

	userResp, profileResp, err := h.authUsecase.Login(ctx, &dto.LoginReq{
		EmailOrUsername: userGoogle.Email,
		Password:        userGoogle.ID,
		RememberMe:      true,
		Oauth2:          true,
	})
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	token, err := h.securityUsecase.JwtRegistredRTAT(ctx, &dto.JwtRegisteredTokenReq{
		AppId:      appID,
		UserId:     userResp.ID,
		RememberMe: true,
	})

	resp := mapper.ResponseSuccess{
		Data: map[string]any{
			"user":    userResp,
			"profile": profileResp,
			"token":   token,
		},
	}

	mapper.NewSuccessResp(w, r, resp, 200)
}
