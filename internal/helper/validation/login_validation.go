package validation

import (
	"fmt"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	_error "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/error"
)

func LoginValidation(req *dto.LoginReq) error {
	err400 := map[string][]string{}

	if req.EmailOrUsername == "" {
		err400["email_or_username"] = append(err400["email_or_username"], fmt.Sprintf(required, "email_or_username"))
	}
	if len(req.EmailOrUsername) < 3 {
		err400["email_or_username"] = append(err400["email_or_username"], fmt.Sprintf(min, "email_or_username", "3"))
	}
	if len(req.EmailOrUsername) > 55 {
		err400["email_or_username"] = append(err400["email_or_username"], fmt.Sprintf(max, "email_or_username", "55"))
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

	if len(err400) != 0 {
		return _error.Err400(err400)
	}
	return nil
}
