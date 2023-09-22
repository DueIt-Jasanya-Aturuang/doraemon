package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jasanya-tech/jasanya-response-backend-golang/_error"
	"github.com/jasanya-tech/jasanya-response-backend-golang/response"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/api/rest/helper"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

type AppMiddleware struct {
	appUsecase domain.AppUsecase
}

func NewAppMiddleware(
	appUsecase domain.AppUsecase,
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

		err := a.appUsecase.CheckByID(r.Context(), &domain.RequestCheckApp{
			AppID: appID,
		})
		if err != nil {
			helper.ErrorResponseEncode(w, _error.HttpErrString("invalid app id", response.CM05))
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
