package validation

import (
	"regexp"

	"github.com/jasanya-tech/jasanya-response-backend-golang/_error"
	"github.com/jasanya-tech/jasanya-response-backend-golang/response"
	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

func GenerateOTPValidation(req *domain.RequestGenerateOTP) error {
	if req.Type != util.ActivasiAccount && req.Type != util.ForgotPassword {
		log.Warn().Msgf("invalid type otp %s", req.Type)
		return _error.HttpErrString("invalid type otp", response.CM05)
	}

	errBadRequest := map[string][]string{}
	if req.Email == "" {
		errBadRequest["email"] = append(errBadRequest["email"], required)
	}
	email := maxMinString(req.Email, 12, 55)
	if email != "" {
		errBadRequest["email"] = append(errBadRequest["email"], email)
	}

	match, err := regexp.MatchString(`^([A-Za-z.]|[0-9])+@gmail.com$`, req.Email)
	if err != nil || !match {
		errBadRequest["email"] = append(errBadRequest["email"], emailMsg)
	}

	if len(errBadRequest) != 0 {
		return _error.HttpErrMapOfSlices(errBadRequest, response.CM06)
	}

	return nil
}

func OTPValidation(req *domain.RequestValidationOTP) error {

	if req.Type != util.ActivasiAccount && req.Type != util.ForgotPassword {
		log.Warn().Msgf("invalid type otp %s", req.Type)
		return _error.HttpErrString("invalid type otp", response.CM05)
	}

	errBadRequest := map[string][]string{}
	if len(req.OTP) != 6 {
		errBadRequest["otp"] = append(errBadRequest["otp"], "kode otp anda tidak valid")
	}

	if req.Email == "" {
		errBadRequest["email"] = append(errBadRequest["email"], required)
	}
	email := maxMinString(req.Email, 12, 55)
	if email != "" {
		errBadRequest["email"] = append(errBadRequest["email"], email)
	}
	match, err := regexp.MatchString(`^([A-Za-z.]|[0-9])+@gmail.com$`, req.Email)
	if err != nil || !match {
		errBadRequest["email"] = append(errBadRequest["email"], emailMsg)
	}

	if len(errBadRequest) != 0 {
		return _error.HttpErrMapOfSlices(errBadRequest, response.CM06)
	}

	return nil
}
