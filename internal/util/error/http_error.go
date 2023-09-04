package _error

import (
	"net/http"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
)

func Err400(msg map[string][]string) error {
	return &model.ErrResponseHTTP{
		Code:    http.StatusBadRequest,
		Message: msg,
	}
}

func ErrString(msg string, code int) error {
	return &model.ErrResponseHTTP{
		Code:    code,
		Message: msg,
	}
}
