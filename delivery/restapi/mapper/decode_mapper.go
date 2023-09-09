package mapper

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"

	_error "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/error"
	_msg "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/msg"
)

func DecodeJson(r *http.Request, data any) error {
	err := json.NewDecoder(r.Body).Decode(data)
	if err == io.EOF {
		log.Err(err).Msg(_msg.LogErrJsonNewDecoderDecode)
		return _error.Err400(map[string][]string{
			"bad_request": {
				"empty body request",
			},
		})
	}

	if err != nil {
		log.Err(err).Msg(_msg.LogErrJsonNewDecoderDecode)
		return err
	}

	return nil

}
