package mapper

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
	_error "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/error"
)

type ResponseError struct {
	Errors *any `json:"errors,omitempty"`
}

// ResponseSuccess http success response
type ResponseSuccess struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func NewErrorResp(w http.ResponseWriter, _ *http.Request, err error) {
	var (
		errHTTP            *model.ErrResponseHTTP
		unmarshalTypeError *json.UnmarshalTypeError
		syntaxError        *json.SyntaxError
	)

	switch {
	case errors.As(err, &unmarshalTypeError):
		err = _error.Err422(map[string][]string{
			err.(*json.UnmarshalTypeError).Field: {
				"invalid type input, type must be %s", err.(*json.UnmarshalTypeError).Type.String(),
			},
		})
	case errors.As(err, &syntaxError):
		err = _error.Err400(map[string][]string{
			"unexpected": {
				"unexpected end of json input",
			},
		})
	case errors.Is(err, context.DeadlineExceeded):
		err = _error.ErrString("Request Timeout", http.StatusRequestTimeout)
	}

	log.Err(err).Msgf("err")
	ok := errors.As(err, &errHTTP)
	if !ok {
		err = _error.ErrStringDefault(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errHTTP.Code)
	resp := ResponseError{
		Errors: &errHTTP.Message,
	}

	if errEncode := json.NewEncoder(w).Encode(resp); errEncode != nil {
		log.Err(errEncode).Msg("failed encode error to json")
	}
}

func NewSuccessResp(w http.ResponseWriter, _ *http.Request, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if errEncode := json.NewEncoder(w).Encode(data); errEncode != nil {
		log.Err(errEncode).Msg("failed encode error to json")
	}
}
