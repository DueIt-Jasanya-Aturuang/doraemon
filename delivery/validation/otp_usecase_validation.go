package validation

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	_error "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/error"
)

func OTPGenerateValidation(req *dto.OTPGenerateReq) error {
	if req.Type != "activasi-account" && req.Type != "forgot-password" {
		log.Warn().Msgf("invalid type otp %s", req.Type)
		return _error.ErrStringDefault(http.StatusForbidden)
	}

	err400 := map[string][]string{}
	if req.Email == "" {
		err400["email"] = append(err400["email"], fmt.Sprintf(required, "email"))
	}
	if len(req.Email) < 12 {
		err400["email"] = append(err400["email"], fmt.Sprintf(min, "email", "12"))
	}
	if len(req.Email) > 55 {
		err400["email"] = append(err400["email"], fmt.Sprintf(max, "email", "55"))
	}
	match, err := regexp.MatchString(`^([A-Za-z.]|[0-9])+@gmail.com$`, req.Email)
	if err != nil || !match {
		err400["email"] = append(err400["email"], email)
	}

	if len(err400) != 0 {
		return _error.Err400(err400)
	}

	return nil
}

func OTPValidation(req *dto.OTPValidationReq) error {

	if req.Type != "activasi-account" && req.Type != "forgot-password" {
		log.Warn().Msgf("invalid type otp %s", req.Type)
		return _error.ErrStringDefault(http.StatusForbidden)
	}

	err400 := map[string][]string{}
	if len(req.OTP) != 6 {
		err400["otp"] = append(err400["otp"], "kode otp anda tidak valid")
	}

	if req.Email == "" {
		err400["email"] = append(err400["email"], fmt.Sprintf(required, "email"))
	}
	if len(req.Email) < 12 {
		err400["email"] = append(err400["email"], fmt.Sprintf(min, "email", "12"))
	}
	if len(req.Email) > 55 {
		err400["email"] = append(err400["email"], fmt.Sprintf(max, "email", "55"))
	}
	match, err := regexp.MatchString(`^([A-Za-z.]|[0-9])+@gmail.com$`, req.Email)
	if err != nil || !match {
		err400["email"] = append(err400["email"], email)
	}

	if len(err400) != 0 {
		return _error.Err400(err400)
	}

	return nil
}
