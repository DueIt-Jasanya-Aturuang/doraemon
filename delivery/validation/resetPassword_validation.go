package validation

import (
	"fmt"
	"net/http"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	_error "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/error"
)

func ResetPasswordValidation(req *dto.ResetPasswordReq) error {
	err400 := map[string][]string{}

	if req.UserID == "" {
		return _error.ErrStringDefault(http.StatusForbidden)
	}
	
	if req.OldPassword == "" {
		err400["old_password"] = append(err400["old_password"], fmt.Sprintf(required, "old_password"))
	}
	if len(req.OldPassword) < 6 {
		err400["old_password"] = append(err400["old_password"], fmt.Sprintf(min, "old_password", "6"))
	}
	if len(req.OldPassword) > 55 {
		err400["old_password"] = append(err400["old_password"], fmt.Sprintf(max, "old_password", "55"))
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
