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
	// set time out proccess
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// set variable request, and header
	var validateAccessReq dto.JwtTokenReq
	appID := r.Header.Get("App-ID")
	userID := r.Header.Get("User-ID")
	token := r.Header.Get("Authorization")

	// check apakah data header atau gak
	// jika tidak ada akan return 401
	if appID == "" || userID == "" || token == "" {
		log.Warn().Msgf("app id / user id / authorization header tidak tersedia")
		mapper.NewErrorResp(w, r, _error.ErrStringDefault(http.StatusUnauthorized))
		return
	}

	// masukan variable header tadi kedalam request validate
	validateAccessReq.AppId = appID
	validateAccessReq.UserId = userID
	validateAccessReq.Authorization = token

	// check appid, jika error akan return error
	// ini error sudah di set dari usecase, apakah error tersebut 500 atau yang lainnya
	err := h.appUsecase.CheckAppByID(ctx, &dto.AppReq{
		AppID: appID,
	})
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	// get url path nya
	// validasi apakah access token valid atau gak
	activasiHeader := r.Header.Get("Activasi")
	expAT, err := h.securityUsecase.JwtValidateAT(ctx, &validateAccessReq, activasiHeader)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	// jika access token expired maka akan registrasi at dan rt ulang
	if expAT {
		// process registrasi rt at
		newAT, err := h.securityUsecase.JwtGenerateRTAT(ctx, &validateAccessReq)
		if err != nil {
			mapper.NewErrorResp(w, r, err)
			return
		}

		log.Info().Msgf("set token baru for user %s", validateAccessReq.UserId)
		// set new token ke dalam header authorization
		w.Header().Set("Authorization", newAT.Token)
	}

	mapper.NewSuccessResp(w, r, nil, 200)
}

func (h *SecurityHandlerImpl) Logout(w http.ResponseWriter, r *http.Request) {
	// set time out proccess
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// declare variable dan get header user id dan authorization
	var reqLogout dto.LogoutReq
	userID := r.Header.Get("User-ID")
	token := r.Header.Get("Authorization")

	// declare response success
	// jika header tersebut kosong maka akan return aja success
	resp := mapper.ResponseSuccess{
		Message: "anda berhasil logout",
	}
	if userID == "" && token == "" {
		log.Warn().Msgf("user id / authorization header tidak tersedia")
		mapper.NewSuccessResp(w, r, resp, 200)
		return
	}

	// set variable header kedalam request
	reqLogout.UserID = userID
	reqLogout.Token = token

	// process logout
	err := h.securityUsecase.Logout(ctx, &reqLogout)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	mapper.NewSuccessResp(w, r, resp, 200)
}
