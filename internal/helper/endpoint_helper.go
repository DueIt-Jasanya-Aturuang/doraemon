package helper

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
)

func EndPointMarshal() string {
	endpoint := []string{
		"/auth/login",
		"/auth/register",
	}

	marshal, err := json.Marshal(endpoint)
	if err != nil {
		log.Err(err).Msg("failed to marshal endpoint")
		return ""
	}

	return string(marshal)
}
