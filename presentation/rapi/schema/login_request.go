package schema

import (
	"github.com/jasanya-tech/jasanya-response-backend-golang/_error"
	"github.com/jasanya-tech/jasanya-response-backend-golang/response"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

type RequestLogin struct {
	EmailOrUsername string `json:"email_or_username"`
	Password        string `json:"password"`
	RememberMe      bool   `json:"remember_me"`
}

func (r *RequestLogin) Validation() error {
	errBadRequest := map[string][]string{}
	if r.EmailOrUsername == "" {
		errBadRequest["email_or_username"] = append(errBadRequest["email_or_username"], util.Required)
	}
	emailOrUsername := util.MaxMinString(r.EmailOrUsername, 3, 55)
	if emailOrUsername != "" {
		errBadRequest["email_or_username"] = append(errBadRequest["email_or_username"], emailOrUsername)
	}

	if r.Password == "" {
		errBadRequest["password"] = append(errBadRequest["password"], util.Required)
	}
	password := util.MaxMinString(r.Password, 6, 55)
	if password != "" {
		errBadRequest["password"] = append(errBadRequest["password"], password)
	}

	if len(errBadRequest) != 0 {
		return _error.HttpErrMapOfSlices(errBadRequest, response.CM06)
	}
	return nil
}
