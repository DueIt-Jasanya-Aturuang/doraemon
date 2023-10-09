package schema

import (
	"github.com/jasanya-tech/jasanya-response-backend-golang/_error"
	"github.com/jasanya-tech/jasanya-response-backend-golang/response"
)

type RequestLoginGoogle struct {
	Token  string `json:"token"`
	Device string `json:"device"`
}

func (r *RequestLoginGoogle) Validation() error {
	errBadRequest := map[string][]string{}
	if r.Token == "" {
		return _error.HttpErrString("invalid token", response.CM05)
	}

	if r.Device != "web" && r.Device != "mobile" {
		return _error.HttpErrString("invalid device", response.CM05)
	}

	if len(errBadRequest) != 0 {
		return _error.HttpErrMapOfSlices(errBadRequest, response.CM06)
	}
	return nil
}
