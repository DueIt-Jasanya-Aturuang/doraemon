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

func Err422(msg map[string][]string) error {
	return &model.ErrResponseHTTP{
		Code:    http.StatusUnprocessableEntity,
		Message: msg,
	}
}

func ErrStringDefault(code int) error {
	var msg string

	switch code {
	case 500:
		msg = "INTERNAL SERVER ERROR"
	case 502:
		msg = "BAD GATEWAY"
	case 404:
		msg = "DATA NOT FOUND"
	case 403:
		msg = "FORBIDDEN"
	case 401:
		msg = "UNAUTHORIZATION"

	}

	return &model.ErrResponseHTTP{
		Code:    code,
		Message: msg,
	}
}

func ErrString(msg string, code int) error {
	if msg == "" {
		switch code {
		case 404:
			msg = "DATA NOT FOUND"
		case 403:
			msg = "FORBIDDEN"
		case 401:
			msg = "UNAUTHORIZATION"
		case 500:
			msg = "INTERNAL SERVER ERROR"
		}
	}

	return &model.ErrResponseHTTP{
		Code:    code,
		Message: msg,
	}
}
