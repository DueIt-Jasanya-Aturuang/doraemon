// Package validation is an example package.
//
// It provides a method to greet someone.
package validation

import (
	"fmt"
	"regexp"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	_error "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/error"
)

/*
RegisterValidation This is a function to validate the dto.RegisterReq and
will be return error for example [_error.Err400]

# Example

	req := *dto.RegisterReq{
		// this is the register request field
		// you can check more details in dto.RegisterReq
	}
	err := RegisterValidation(req)
	if err != nil {
		// this will return the return of package _error
		// you can check more details in _error(http_error.go)
	}
*/
func RegisterValidation(req *dto.RegisterReq) error {
	err400 := map[string][]string{}

	if req.FullName == "" {
		err400["full_name"] = append(err400["full_name"], fmt.Sprintf(required, "full_name"))
	}
	if len(req.FullName) < 3 {
		err400["full_name"] = append(err400["full_name"], fmt.Sprintf(min, "full_name", "3"))
	}
	if len(req.FullName) > 32 {
		err400["full_name"] = append(err400["full_name"], fmt.Sprintf(max, "full_name", "32"))
	}

	if req.Username == "" {
		err400["username"] = append(err400["username"], fmt.Sprintf(required, "username"))
	}
	if len(req.Username) < 3 {
		err400["username"] = append(err400["username"], fmt.Sprintf(min, "username", "3"))
	}
	if len(req.Username) > 22 {
		err400["username"] = append(err400["username"], fmt.Sprintf(max, "username", "22"))
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

	if len(err400) != 0 {
		return _error.Err400(err400)
	}
	return nil

}
