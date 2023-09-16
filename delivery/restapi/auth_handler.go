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

type AuthHandlerImpl struct {
	authUsecase usecase.AuthUsecase
	appUsecase  usecase.AppUsecase
	otpUsecase  usecase.OTPUsecase
}

func NewAuthHandlerImpl(
	authUsecase usecase.AuthUsecase,
	appUsecase usecase.AppUsecase,
	otpUsecase usecase.OTPUsecase,
) *AuthHandlerImpl {
	return &AuthHandlerImpl{
		authUsecase: authUsecase,
		appUsecase:  appUsecase,
		otpUsecase:  otpUsecase,
	}
}

func (h *AuthHandlerImpl) Register(w http.ResponseWriter, r *http.Request) {
	// set time out proccess
	// testing in 4.776681149s
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
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
	var reqRegister dto.RegisterReq
	err = mapper.DecodeJson(r, &reqRegister)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	// set role, appid, dan email verified
	reqRegister.Role = 1
	reqRegister.AppID = appID
	reqRegister.EmailVerifiedAt = false

	// validasi request
	err = validation.RegisterValidation(&reqRegister)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	// register user
	userResp, profileResp, token, err := h.authUsecase.Register(ctx, &reqRegister)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	// mencoba untuk push otp langsung, ga peduli error atau tidak
	_ = h.otpUsecase.OTPGenerate(ctx, &dto.OTPGenerateReq{
		Email:  reqRegister.Email,
		Type:   "activasi-account",
		UserID: userResp.ID,
	})

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
	// set time out proccess
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
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
	var reqLogin dto.LoginReq
	err = mapper.DecodeJson(r, &reqLogin)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}
	reqLogin.AppID = appID
	// validasi request
	err = validation.LoginValidation(&reqLogin)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	// login user
	userResp, profileResp, token, err := h.authUsecase.Login(ctx, &reqLogin)
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
