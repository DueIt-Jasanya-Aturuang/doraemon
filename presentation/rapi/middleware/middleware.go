package middleware

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jasanya-tech/jasanya-response-backend-golang/_error"
	"github.com/jasanya-tech/jasanya-response-backend-golang/response"
	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/presentation/rapi/helper"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

type AppMiddleware struct {
	appUsecase usecase.AppUsecase
}

func NewAppMiddleware(
	appUsecase usecase.AppUsecase,
) *AppMiddleware {
	return &AppMiddleware{
		appUsecase: appUsecase,
	}
}

func (a *AppMiddleware) CheckAppID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appID := r.Header.Get(util.AppIDHeader)
		if _, err := uuid.Parse(appID); err != nil {
			helper.ErrorResponseEncode(w, _error.HttpErrString("invalid app id", response.CM05))
			return
		}

		err := a.appUsecase.CheckByID(r.Context(), appID)
		if err != nil {
			if errors.Is(err, usecase.InvalidAppID) {
				err = _error.HttpErrString("invalid app id", response.CM05)
			}
			helper.ErrorResponseEncode(w, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func SetAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newToken := r.Header.Get("Authorization")
		if newToken != "" {
			w.Header().Set("Authorization", newToken)
		}
		next.ServeHTTP(w, r)
	})
}

func CheckApiKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-Api-Key")
		if apiKey != infra.AppApiKeyAuth {
			log.Warn().Msgf("invalid old key | key : %s", apiKey)
			helper.ErrorResponseEncode(w, _error.HttpErrString("forbidden", response.CM05))
			return
		}
		next.ServeHTTP(w, r)
	})
}
