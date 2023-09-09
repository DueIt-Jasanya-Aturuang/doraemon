package restapi

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/delivery/restapi/mapper"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/delivery/restapi/middleware"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/delivery/validation"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/usecase"
	_error "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/error"
)

type OTPHandlerImpl struct {
	otpUsecase usecase.OTPUsecase
	appUsecase usecase.AppUsecase
}

func NewOTPHandlerImpl(
	otpUsecase usecase.OTPUsecase,
	appUsecase usecase.AppUsecase,
) *OTPHandlerImpl {
	return &OTPHandlerImpl{
		otpUsecase: otpUsecase,
		appUsecase: appUsecase,
	}
}

func (h *OTPHandlerImpl) OTPGenerate(w http.ResponseWriter, r *http.Request) {
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
	// and set type dari header type
	var reqOTPGenerate dto.OTPGenerateReq
	err = mapper.DecodeJson(r, &reqOTPGenerate)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}
	reqOTPGenerate.Type = r.Header.Get("Type")

	// validasi request
	err = validation.OTPGenerateValidation(&reqOTPGenerate)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	// check apakan type nya activasi atau tidak
	// jika activasi maka akan get userid di header, jika tidak ada maka akan return fobidden
	if reqOTPGenerate.Type == "activasi-account" {
		userID := r.Header.Get("User-ID")
		if userID == "" {
			mapper.NewErrorResp(w, r, _error.ErrStringDefault(http.StatusForbidden))
			return
		}

		reqOTPGenerate.UserID = userID
	}

	// generate otp
	err = h.otpUsecase.OTPGenerate(ctx, &reqOTPGenerate)
	if err != nil {
		middleware.DeletedClientHelper(reqOTPGenerate.Email + ":" + reqOTPGenerate.Type)
		mapper.NewErrorResp(w, r, err)
		return
	}

	// set limiter
	err = middleware.RateLimiterOTP(&reqOTPGenerate)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	resp := mapper.ResponseSuccess{
		Message: "kode otp telah berhasil dikirim, silahkan cek gmail anda",
	}

	mapper.NewSuccessResp(w, r, resp, 200)
}
