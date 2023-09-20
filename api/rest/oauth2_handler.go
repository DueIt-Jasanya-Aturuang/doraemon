package rest

import (
	"errors"
	"net/http"

	"github.com/jasanya-tech/jasanya-response-backend-golang/_error"
	"github.com/jasanya-tech/jasanya-response-backend-golang/response"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/api/rest/helper"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/api/validation"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/_usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

type Oauth2HandlerImpl struct {
	oauth2Usecase domain.Oauth2Usecase
	authUsecase   domain.AuthUsecase
}

func NewOauth2HandlerImpl(
	oauth2Usecase domain.Oauth2Usecase,
	authUsecase domain.AuthUsecase,
) *Oauth2HandlerImpl {
	return &Oauth2HandlerImpl{
		oauth2Usecase: oauth2Usecase,
		authUsecase:   authUsecase,
	}
}

func (h *Oauth2HandlerImpl) LoginWithGoogle(w http.ResponseWriter, r *http.Request) {
	req := new(domain.RequestLoginWithGoogle)
	err := helper.DecodeJson(r, req)
	if err != nil {
		helper.ErrorResponseEncode(w, err)
		return
	}

	err = validation.Oauth2LoginWithGoogleValidation(req)
	if err != nil {
		helper.ErrorResponseEncode(w, err)
		return
	}

	userGoogle, err := h.oauth2Usecase.GoogleClaimUser(r.Context(), req)
	if err != nil {
		if errors.Is(err, _usecase.InvalidTokenOauth) {
			err = _error.HttpErrString("invalid token", response.CM05)
		}
		helper.ErrorResponseEncode(w, err)
		return
	}

	// jika user tidak ada maka akan register -> login -> return
	appID := r.Header.Get(util.AppIDHeader)
	if !userGoogle.ExistsUser {
		// register request
		reqRegister := &domain.RequestRegister{
			FullName:        userGoogle.Name,
			Username:        userGoogle.GivenName,
			Email:           userGoogle.Email,
			Password:        userGoogle.ID,
			RePassword:      userGoogle.ID,
			EmailVerifiedAt: true,
			AppID:           appID,
			Role:            1,
		}

		// register process
		resp, err := h.authUsecase.Register(r.Context(), reqRegister)
		if err != nil {
			if errors.Is(err, _usecase.EmailIsExist) {
				err = _error.HttpErrMapOfSlices(map[string][]string{
					"email": {
						"email sudah terdaftar",
					},
				}, response.CM06)
			}
			if errors.Is(err, _usecase.UsernameIsExist) {
				err = _error.HttpErrMapOfSlices(map[string][]string{
					"username": {
						"username sudah tersedia",
					},
				}, response.CM06)
			}
			helper.ErrorResponseEncode(w, err)
			return
		}

		helper.SuccessResponseEncode(w, resp, "login successfully")
		return
	}

	// jika ada maka langsung login
	resp, err := h.authUsecase.Login(r.Context(), &domain.RequestLogin{
		EmailOrUsername: userGoogle.Email,
		Password:        userGoogle.ID,
		RememberMe:      true,
		Oauth2:          true,
	})
	if err != nil {
		helper.ErrorResponseEncode(w, err)
		return
	}

	helper.SuccessResponseEncode(w, resp, "login successfully")
}
