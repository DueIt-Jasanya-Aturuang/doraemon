package validation

import (
	"github.com/jasanya-tech/jasanya-response-backend-golang/_error"
	"github.com/jasanya-tech/jasanya-response-backend-golang/response"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain"
)

func Oauth2LoginWithGoogleValidation(req *domain.RequestLoginWithGoogle) error {

	if req.Token == "" {
		return _error.HttpErrString("invalid token", response.CM05)
	}

	if req.Device != "web" && req.Device != "mobile" {
		return _error.HttpErrString("invalid device", response.CM05)
	}

	return nil
}
