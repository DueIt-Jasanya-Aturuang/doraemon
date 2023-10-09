package helper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jasanya-tech/jasanya-response-backend-golang/_error"
	"github.com/jasanya-tech/jasanya-response-backend-golang/response"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

func ErrorResponseEncode(w http.ResponseWriter, err error) {
	var (
		errHttp          *response.HttpError
		errUnmarshalType *json.UnmarshalTypeError
		errSyntak        *json.SyntaxError
		errPQ            *pq.Error
	)

	switch {
	case errors.As(err, &errUnmarshalType):
		msg := fmt.Sprintf("UnprocessableEntity : %s", err.Error())
		err = _error.HttpErrString(msg, response.CM11)

	case errors.As(err, &errSyntak):
		msg := fmt.Sprintf("unexpected end of json input : %s", err.Error())
		err = _error.HttpErrString(msg, response.CM11)

	case errors.As(err, &errPQ):
		log.Warn().Msgf("pqerror | err : %v", err)
		if errPQ.Code == "23503" {
			err = _error.HttpErrString(response.CodeCompanyName[response.CM05], response.CM05)
		} else {
			err = _error.HttpErrString(response.CodeCompanyName[response.CM99], response.CM99)
		}
	case errors.Is(err, context.DeadlineExceeded):
		err = _error.HttpErrString(response.CodeCompanyName[response.CM08], response.CM08)
	}

	ok := errors.As(err, &errHttp)
	if !ok {
		err = _error.HttpErrString(string(response.CM99), response.CM99)
		errors.As(err, &errHttp)
	}

	resp := response.Error(*errHttp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Code)

	if errEncode := json.NewEncoder(w).Encode(resp); errEncode != nil {
		log.Err(errEncode).Msgf(util.LogErrEncode, resp, errEncode)
	}
}

func SuccessResponseEncode(w http.ResponseWriter, data any, message string) {
	resp := response.Success(data, response.CM00, message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Code)

	if errEncode := json.NewEncoder(w).Encode(resp); errEncode != nil {
		log.Err(errEncode).Msgf(util.LogErrEncode, resp, errEncode)
	}
}
