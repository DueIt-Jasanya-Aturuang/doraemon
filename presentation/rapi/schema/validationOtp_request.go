package schema

import (
	"regexp"

	"github.com/jasanya-tech/jasanya-response-backend-golang/_error"
	"github.com/jasanya-tech/jasanya-response-backend-golang/response"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

type RequestValidationOTP struct {
	Email string `json:"email"`
	OTP   string `json:"otp_usecase"`
}

func (r *RequestValidationOTP) Validation() error {
	// if req.Type != util.ActivasiAccount && req.Type != util.ForgotPassword {
	// 	log.Warn().Msgf("invalid type otp_usecase %s", req.Type)
	// 	return _error.HttpErrString("invalid type otp_usecase", response.CM05)
	// }

	errBadRequest := map[string][]string{}
	if len(r.OTP) != 6 {
		errBadRequest["otp_usecase"] = append(errBadRequest["otp_usecase"], "kode otp_usecase anda tidak valid")
	}

	if r.Email == "" {
		errBadRequest["email"] = append(errBadRequest["email"], util.Required)
	}
	email := util.MaxMinString(r.Email, 12, 55)
	if email != "" {
		errBadRequest["email"] = append(errBadRequest["email"], email)
	}
	match, err := regexp.MatchString(`^([A-Za-z.]|[0-9])+@gmail.com$`, r.Email)
	if err != nil || !match {
		errBadRequest["email"] = append(errBadRequest["email"], util.EmailMsg)
	}

	if len(errBadRequest) != 0 {
		return _error.HttpErrMapOfSlices(errBadRequest, response.CM06)
	}

	return nil
}
