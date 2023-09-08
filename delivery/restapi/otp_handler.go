package restapi

import (
	"context"
	"net/http"
	"time"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/delivery/restapi/mapper"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/delivery/restapi/middleware"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/delivery/validation"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/usecase"
	_error "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/error"
)

type OTPHandlerImpl struct {
	otpUsecase usecase.OTPUsecase
}

func NewOTPHandlerImpl(
	otpUsecase usecase.OTPUsecase,
) *OTPHandlerImpl {
	return &OTPHandlerImpl{
		otpUsecase: otpUsecase,
	}
}

func (h *OTPHandlerImpl) OTPGenerate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var reqOTPGenerate dto.OTPGenerateReq
	err := mapper.DecodeJson(r, &reqOTPGenerate)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}
	reqOTPGenerate.Type = r.Header.Get("Type")

	err = validation.OTPGenerateValidation(&reqOTPGenerate)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	if reqOTPGenerate.Type == "activasi-account" {
		userID := r.Header.Get("User-ID")
		if userID == "" {
			mapper.NewErrorResp(w, r, _error.ErrStringDefault(http.StatusForbidden))
			return
		}

		reqOTPGenerate.UserID = userID
	}

	err = middleware.RateLimiterOTP(&reqOTPGenerate)
	if err != nil {
		mapper.NewErrorResp(w, r, err)
		return
	}

	err = h.otpUsecase.OTPGenerate(ctx, &reqOTPGenerate)
	if err != nil {
		middleware.DeletedClientHelper(reqOTPGenerate.Email + ":" + reqOTPGenerate.Type)
		mapper.NewErrorResp(w, r, err)
		return
	}

	resp := mapper.ResponseSuccess{
		Message: "kode otp telah berhasil dikirim, silahkan cek gmail anda",
	}

	mapper.NewSuccessResp(w, r, resp, 200)
}
