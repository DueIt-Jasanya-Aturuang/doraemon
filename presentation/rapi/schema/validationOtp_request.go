package schema

import (
	"regexp"

	"github.com/jasanya-tech/jasanya-response-backend-golang/_error"
	"github.com/jasanya-tech/jasanya-response-backend-golang/response"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

type RequestValidationOTP struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

func (r *RequestValidationOTP) Validation() error {
	errBadRequest := map[string][]string{}
	if len(r.OTP) != 6 {
		errBadRequest["otp"] = append(errBadRequest["otp"], "kode otp anda tidak valid")
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
