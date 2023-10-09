package rapi

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/jasanya-tech/jasanya-response-backend-golang/_error"
	"github.com/jasanya-tech/jasanya-response-backend-golang/response"
	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/presentation/rapi/helper"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

func (p *Presenter) ValidateAccess(w http.ResponseWriter, r *http.Request) {
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
		log.Info().Msgf("invalid activasi header | err : %v", err)
		helper.ErrorResponseEncode(w, _error.HttpErrString(response.CodeCompanyName[response.CM05], response.CM05))
	}

	expAT, err := p.securityUsecase.ValidationJWT(r.Context(), &usecase.RequestValidationJWT{
		AppID:          appID,
		Authorization:  token,
		UserID:         userID,
		ActivasiHeader: activasiHeaderBool,
	})
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

	if expAT {
		newAT, err := p.securityUsecase.ReGenerateJWT(r.Context(), &usecase.RequestValidationJWT{
			AppID:          appID,
			Authorization:  token,
			UserID:         userID,
			ActivasiHeader: activasiHeaderBool,
		})
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

		log.Info().Msgf("set token baru for user %s", userID)

		w.Header().Set("Authorization", newAT.Token)
	} else {
		w.Header().Set("Authorization", r.Header.Get("Authorization"))
	}

	helper.SuccessResponseEncode(w, nil, "ok")
}

func (p *Presenter) Logout(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get(util.UserIDHeader)
	token := r.Header.Get(util.AuthorizationHeader)

	if userID == "" && token == "" {
		log.Warn().Msgf("user id / authorization header tidak tersedia")
		helper.SuccessResponseEncode(w, nil, "successfully logouts")
		return
	}

	err := p.securityUsecase.Logout(r.Context(), &usecase.RequestLogout{
		Token:  token,
		UserID: userID,
	})
	if err != nil {
		helper.ErrorResponseEncode(w, err)
		return
	}

	helper.SuccessResponseEncode(w, nil, "successfully logouts")
}
