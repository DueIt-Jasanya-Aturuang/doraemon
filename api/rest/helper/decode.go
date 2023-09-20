package helper

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/jasanya-tech/jasanya-response-backend-golang/_error"
	"github.com/jasanya-tech/jasanya-response-backend-golang/response"
	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

func DecodeJson(r *http.Request, data any) error {
	err := json.NewDecoder(r.Body).Decode(data)
	if err == io.EOF {
		return _error.HttpErrMapOfSlices(map[string][]string{
			"bad_request": {
				"tidak ada request body",
			},
		}, response.CM06)
	}

	if err != nil {
		log.Warn().Msgf(util.LogErrDecode, r.Body, err)
		return err
	}

	return nil

}
