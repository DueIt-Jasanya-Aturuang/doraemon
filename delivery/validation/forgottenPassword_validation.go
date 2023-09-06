package validation

import (
	"fmt"
	"regexp"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	_error "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/error"
)

func ResetForgottenPasswordValidation(req *dto.ResetForgottenPasswordReq) error {
	err400 := map[string][]string{}

	if req.Token == "" {
		err400["token"] = append(err400["token"], fmt.Sprintf(required, "token"))
	}

	if req.Password == "" {
		err400["password"] = append(err400["password"], fmt.Sprintf(required, "password"))
	}
	if len(req.Password) < 6 {
		err400["password"] = append(err400["password"], fmt.Sprintf(min, "password", "6"))
	}
	if len(req.Password) > 55 {
		err400["password"] = append(err400["password"], fmt.Sprintf(max, "password", "55"))
	}
	if req.Password != req.RePassword {
		err400["password"] = append(err400["password"], passwordAndRePassword)
		err400["re_password"] = append(err400["re_password"], passwordAndRePassword)
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

func ForgottenPasswordValidation(req *dto.ForgottenPasswordReq) error {
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
