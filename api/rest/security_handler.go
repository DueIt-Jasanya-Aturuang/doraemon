package rest

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/jasanya-tech/jasanya-response-backend-golang/_error"
	"github.com/jasanya-tech/jasanya-response-backend-golang/response"
	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/api/rest/helper"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

type SecurityHandlerImpl struct {
	securityUsecase domain.SecurityUsecase
}

func NewSecurityHandlerImpl(
	securityUsecase domain.SecurityUsecase,
) *SecurityHandlerImpl {
	return &SecurityHandlerImpl{
		securityUsecase: securityUsecase,
	}
}

func (h *SecurityHandlerImpl) ValidateAccess(w http.ResponseWriter, r *http.Request) {
	req := new(domain.RequestJwtToken)
	appID := r.Header.Get(util.AppIDHeader)
	userID := r.Header.Get(util.UserIDHeader)
	token := r.Header.Get(util.AuthorizationHeader)

	if appID == "" || userID == "" || token == "" {
		log.Warn().Msgf("app id / user id / authorization header tidak tersedia")
		helper.ErrorResponseEncode(w, _error.HttpErrString(response.CodeCompanyName[response.CM04], response.CM04))
		return
	}

	activasiHeader := r.Header.Get(util.ActivasiHeader)
	activasiHeaderBool, err := strconv.ParseBool(activasiHeader)
	if err != nil {
		helper.ErrorResponseEncode(w, _error.HttpErrString(response.CodeCompanyName[response.CM05], response.CM05))
	}

	req.AppId = appID
	req.UserId = userID
	req.Authorization = token
	req.ActivasiHeader = activasiHeaderBool

	expAT, err := h.securityUsecase.JwtValidation(r.Context(), req)
	if err != nil {
		if errors.Is(err, usecase.InvalidToken) {
			err = _error.HttpErrString("invalid token", response.CM04)
		}
		if errors.Is(err, usecase.JwtUserIDAndHeaderUserIDNotMatch) {
			err = _error.HttpErrString("user id tidak sesuai", response.CM04)
		}
		if errors.Is(err, usecase.InvalidUserID) {
			err = _error.HttpErrString("user id tidak valid", response.CM04)
		}
		if errors.Is(err, usecase.UserIsNotActivited) {
			err = _error.HttpErrString("akun anda belum di aktivasi", response.CM05)
		}
		if errors.Is(err, usecase.JwtAppIDAndHeaderAppIDNotMatch) {
			err = _error.HttpErrString("app id tidak sesuai", response.CM05)
		}
		helper.ErrorResponseEncode(w, err)
		return
	}

	log.Info().Msgf("%v", expAT)
	// jika access token expired maka akan registrasi at dan rt ulang
	if expAT {
		newAT, err := h.securityUsecase.JwtGenerate(r.Context(), req)
		if err != nil {
			if errors.Is(err, usecase.InvalidToken) {
				err = _error.HttpErrString("invalid token", response.CM04)
			}
			if errors.Is(err, usecase.JwtUserIDAndHeaderUserIDNotMatch) {
				err = _error.HttpErrString("user id tidak sesuai", response.CM04)
			}
			helper.ErrorResponseEncode(w, err)
			return
		}

		log.Info().Msgf("set token baru for user %s", req.UserId)
		// set new token ke dalam header authorization
		w.Header().Set("Authorization", newAT.Token)
	} else {
		w.Header().Set("Authorization", r.Header.Get("Authorization"))
	}

	helper.SuccessResponseEncode(w, nil, "ok")
}

func (h *SecurityHandlerImpl) Logout(w http.ResponseWriter, r *http.Request) {
	req := new(domain.RequestLogout)
	userID := r.Header.Get(util.UserIDHeader)
	token := r.Header.Get(util.AuthorizationHeader)

	if userID == "" && token == "" {
		log.Warn().Msgf("user id / authorization header tidak tersedia")
		helper.SuccessResponseEncode(w, nil, "successfully logouts")
		return
	}

	// set variable header kedalam request
	req.UserID = userID
	req.Token = token

	// process logout
	err := h.securityUsecase.Logout(r.Context(), req)
	if err != nil {
		helper.ErrorResponseEncode(w, err)
		return
	}

	helper.SuccessResponseEncode(w, nil, "successfully logouts")
}
