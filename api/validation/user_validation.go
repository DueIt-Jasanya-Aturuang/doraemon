package validation

import (
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/jasanya-tech/jasanya-response-backend-golang/_error"
	"github.com/jasanya-tech/jasanya-response-backend-golang/response"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain"
)

func ResetForgottenPasswordValidation(req *domain.RequestResetForgottenPassword) error {
	errBadRequest := map[string][]string{}

	if req.Token == "" {
		errBadRequest["token"] = append(errBadRequest["token"], required)
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

func ResetPasswordValidation(req *domain.RequestChangePassword) error {
	if _, err := uuid.Parse(req.UserID); err != nil {
		return _error.HttpErrString(response.CodeCompanyName[response.CM04], response.CM04)
	}

	errBadRequest := map[string][]string{}

	if req.OldPassword == "" {
		errBadRequest["old_password"] = append(errBadRequest["old_password"], required)
	}
	oldPassword := maxMinString(req.OldPassword, 6, 55)
	if oldPassword != "" {
		errBadRequest["old_password"] = append(errBadRequest["old_password"], oldPassword)
	}

	if req.Password == "" {
		errBadRequest["password"] = append(errBadRequest["password"], fmt.Sprintf(required, "password"))
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
