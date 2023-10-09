package schema

import (
	"github.com/jasanya-tech/jasanya-response-backend-golang/_error"
	"github.com/jasanya-tech/jasanya-response-backend-golang/response"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

type RequestResetForgottenPassword struct {
	Password   string `json:"password"`
	RePassword string `json:"re_password"`
}

func (r *RequestResetForgottenPassword) Validation() error {
	errBadRequest := map[string][]string{}
	if r.Password == "" {
		errBadRequest["password"] = append(errBadRequest["password"], util.Required)
	}
	password := util.MaxMinString(r.Password, 6, 55)
	if password != "" {
		errBadRequest["password"] = append(errBadRequest["password"], password)
	}
	if r.Password != r.RePassword {
		errBadRequest["password"] = append(errBadRequest["password"], util.PasswordAndRePassword)
		errBadRequest["re_password"] = append(errBadRequest["re_password"], util.PasswordAndRePassword)
	}

	if len(errBadRequest) != 0 {
		return _error.HttpErrMapOfSlices(errBadRequest, response.CM06)
	}
	return nil
}
