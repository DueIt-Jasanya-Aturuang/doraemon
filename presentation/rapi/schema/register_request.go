package schema

import (
	"regexp"

	"github.com/jasanya-tech/jasanya-response-backend-golang/_error"
	"github.com/jasanya-tech/jasanya-response-backend-golang/response"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

type RequestRegister struct {
	FullName   string `json:"full_name"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	RePassword string `json:"re_password"`
}

func (r *RequestRegister) Validation() error {
	errBadRequest := map[string][]string{}
	if r.FullName == "" {
		errBadRequest["full_name"] = append(errBadRequest["full_name"], util.Required)
	}

	fullName := util.MaxMinString(r.FullName, 3, 32)
	if fullName != "" {
		errBadRequest["full_name"] = append(errBadRequest["full_name"], fullName)
	}

	if r.Username == "" {
		errBadRequest["username"] = append(errBadRequest["username"], util.Required)
	}

	username := util.MaxMinString(r.Username, 3, 22)
	if username != "" {
		errBadRequest["username"] = append(errBadRequest["username"], username)
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
