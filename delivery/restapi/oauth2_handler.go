package restapi

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/delivery/restapi/mapper"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/delivery/validation"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/usecase"
	_error "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/error"
)

type Oauth2HandlerImpl struct {
	oauth2Usecase usecase.Oauth2Usecase
	authUsecase   usecase.AuthUsecase
	appUsecase    usecase.AppUsecase
}

func NewOauth2HandlerImpl(
	oauth2Usecase usecase.Oauth2Usecase,
	authUsecase usecase.AuthUsecase,
	appUsecase usecase.AppUsecase,
) *Oauth2HandlerImpl {
	return &Oauth2HandlerImpl{
		oauth2Usecase: oauth2Usecase,
		authUsecase:   authUsecase,
		appUsecase:    appUsecase,
	}
}

func (h *Oauth2HandlerImpl) LoginWithGoogle(w http.ResponseWriter, r *http.Request) {
	// set time out proccess
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// get header App-ID
	// jika gak ada akan return forbidden
	appID := r.Header.Get("App-ID")
	if appID == "" {
		log.Warn().Msg("tidak ada header appid")
		mapper.NewErrorResp(w, r, _error.ErrStringDefault(http.StatusForbidden))
		return
	}

	// check appid, jika error akan return error
	// ini error sudah di set dari usecase, apakah error tersebut 500 atau yang lainnya
	err := h.appUsecase.CheckAppByID(ctx, &dto.AppReq{
		AppID: appID,
	})
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	// decod request ke dalam dto
	var reqLogin dto.LoginGoogleReq
	err = mapper.DecodeJson(r, &reqLogin)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	// validasi request
	err = validation.Oauth2LoginValidation(&reqLogin)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	// claim google user
	userGoogle, err := h.oauth2Usecase.GoogleClaimUser(ctx, &reqLogin)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	// jika user tidak ada maka akan register -> login -> return
	if !userGoogle.ExistsUser {
		// register request
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

		// register process
		userResp, profileResp, token, err := h.authUsecase.Register(ctx, reqRegister)
		if err != nil {
			mapper.NewErrorResp(w, r, err)
			return
		}

		// response
		resp := mapper.ResponseSuccess{
			Data: map[string]any{
				"user":    userResp,
				"profile": profileResp,
				"token":   token,
			},
		}

		mapper.NewSuccessResp(w, r, resp, 200)
		return
	}

	// jika ada maka langsung login
	userResp, profileResp, token, err := h.authUsecase.Login(ctx, &dto.LoginReq{
		EmailOrUsername: userGoogle.Email,
		Password:        userGoogle.ID,
		RememberMe:      true,
		Oauth2:          true,
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
