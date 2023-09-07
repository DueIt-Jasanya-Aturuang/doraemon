package validation

import (
	"net/http"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	_error "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/error"
)

func Oauth2LoginValidation(req *dto.LoginGoogleReq) error {

	if req.Token == "" {
		return _error.ErrStringDefault(http.StatusForbidden)
	}

	if req.Device != "web" && req.Device != "mobile" {
		return _error.ErrStringDefault(http.StatusForbidden)
	}

	return nil
}
