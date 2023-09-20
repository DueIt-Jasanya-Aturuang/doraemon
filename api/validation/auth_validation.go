package validation

import (
	"regexp"

	"github.com/jasanya-tech/jasanya-response-backend-golang/_error"
	"github.com/jasanya-tech/jasanya-response-backend-golang/response"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain"
)

func LoginValidation(req *domain.RequestLogin) error {
	errBadRequest := map[string][]string{}
	if req.EmailOrUsername == "" {
		errBadRequest["email_or_username"] = append(errBadRequest["email_or_username"], required)
	}
	emailOrUsername := maxMinString(req.EmailOrUsername, 3, 55)
	if emailOrUsername != "" {
		errBadRequest["email_or_username"] = append(errBadRequest["email_or_username"], emailOrUsername)
	}

	if req.Password == "" {
		errBadRequest["password"] = append(errBadRequest["password"], required)
	}
	password := maxMinString(req.Password, 6, 55)
	if password != "" {
		errBadRequest["password"] = append(errBadRequest["password"], password)
	}

	if len(errBadRequest) != 0 {
		return _error.HttpErrMapOfSlices(errBadRequest, response.CM06)
	}
	return nil
}

func RegisterValidation(req *domain.RequestRegister) error {
	errBadRequest := map[string][]string{}
	if req.FullName == "" {
		errBadRequest["full_name"] = append(errBadRequest["full_name"], required)
	}

	fullName := maxMinString(req.FullName, 3, 32)
	if fullName != "" {
		errBadRequest["full_name"] = append(errBadRequest["full_name"], fullName)
	}

	if req.Username == "" {
		errBadRequest["username"] = append(errBadRequest["username"], required)
	}

	username := maxMinString(req.Username, 3, 22)
	if username != "" {
		errBadRequest["username"] = append(errBadRequest["username"], username)
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

	if req.Password == "" {
		errBadRequest["password"] = append(errBadRequest["password"], required)
	}
	password := maxMinString(req.Password, 6, 55)
	if password != "" {
		errBadRequest["password"] = append(errBadRequest["password"], password)
	}
	if req.Password != req.RePassword {
		errBadRequest["password"] = append(errBadRequest["password"], passwordAndRePassword)
		errBadRequest["re_password"] = append(errBadRequest["re_password"], passwordAndRePassword)
	}

	if len(errBadRequest) != 0 {
		return _error.HttpErrMapOfSlices(errBadRequest, response.CM06)
	}
	return nil
}
