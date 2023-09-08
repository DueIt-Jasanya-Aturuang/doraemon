package restapi

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/delivery/restapi/mapper"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/usecase"
	_error "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/error"
)

type SecurityHandlerImpl struct {
	securityUsecase usecase.SecurityUsecase
	appUsecase      usecase.AppUsecase
}

func NewSecurityHandlerImpl(
	securityUsecase usecase.SecurityUsecase,
	appUsecase usecase.AppUsecase,
) *SecurityHandlerImpl {
	return &SecurityHandlerImpl{
		securityUsecase: securityUsecase,
		appUsecase:      appUsecase,
	}
}

func (h *SecurityHandlerImpl) ValidateAccess(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var validateAccessReq dto.JwtTokenReq
	appID := r.Header.Get("App-ID")
	userID := r.Header.Get("User-ID")
	token := r.Header.Get("Authorization")

	if appID == "" || userID == "" || token == "" {
		log.Warn().Msgf("app id / user id / authorization header tidak tersedia")
		mapper.NewErrorResp(w, r, _error.ErrStringDefault(http.StatusUnauthorized))
		return
	}

	validateAccessReq.AppId = appID
	validateAccessReq.UserId = userID
	validateAccessReq.Authorization = token

	err := h.appUsecase.CheckAppByID(ctx, &dto.AppReq{
		AppID: appID,
	})
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	path := r.URL.Path
	expAT, err := h.securityUsecase.JwtValidateAT(ctx, &validateAccessReq, path)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	var newAT *dto.JwtTokenResp
	if expAT {
		newAT, err = h.securityUsecase.JwtGenerateRTAT(ctx, &validateAccessReq)
		if err != nil {
			mapper.NewErrorResp(w, r, err)
			return
		}
	}

	if newAT != nil {
		log.Info().Msgf("set token baru for user %s", validateAccessReq.UserId)
		w.Header().Set("Authorization", newAT.Token)
	}

	mapper.NewSuccessResp(w, r, nil, 200)
}

func (h *SecurityHandlerImpl) Logout(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var reqLogout dto.LogoutReq
	userID := r.Header.Get("User-ID")
	token := r.Header.Get("Authorization")

	resp := mapper.ResponseSuccess{
		Message: "anda berhasil logout",
	}
	if userID == "" || token == "" {
		log.Warn().Msgf("user id / authorization header tidak tersedia")
		mapper.NewSuccessResp(w, r, resp, 200)
		return
	}

	reqLogout.UserID = userID
	reqLogout.Token = token

	err := h.securityUsecase.Logout(ctx, &reqLogout)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	mapper.NewSuccessResp(w, r, resp, 200)
}
