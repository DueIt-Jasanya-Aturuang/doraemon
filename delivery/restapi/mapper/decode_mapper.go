package mapper

import (
	"encoding/json"
	"io"
	"net/http"

	_error "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/error"
)

func DecodeJson(r *http.Request, data any) error {
	err := json.NewDecoder(r.Body).Decode(data)
	if err == io.EOF {
		return _error.Err400(map[string][]string{
			"bad_request": {
				"empty body request",
			},
		})
	}

	if err != nil {
		return err
	}

	return nil

}
