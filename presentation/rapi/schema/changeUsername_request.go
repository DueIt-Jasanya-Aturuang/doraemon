package schema

import (
	"github.com/jasanya-tech/jasanya-response-backend-golang/_error"
	"github.com/jasanya-tech/jasanya-response-backend-golang/response"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

type RequestChangeUsername struct {
	Username string `json:"username"`
}

func (r *RequestChangeUsername) Validation() error {
	errBadRequest := map[string][]string{}

	if r.Username == "" {
		errBadRequest["username"] = append(errBadRequest["username"], util.Required)
	}

	username := util.MaxMinString(r.Username, 3, 22)
	if username != "" {
		errBadRequest["username"] = append(errBadRequest["username"], username)
	}

	if len(errBadRequest) != 0 {
		return _error.HttpErrMapOfSlices(errBadRequest, response.CM06)
	}
	return nil
}
