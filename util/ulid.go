package util

import (
	"github.com/jasanya-tech/jasanya-response-backend-golang/_error"
	"github.com/jasanya-tech/jasanya-response-backend-golang/response"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

func ParseUlid(s string) error {
	if _, err := ulid.Parse(s); err != nil {
		log.Info().Msgf("failed parse ulid | err : %v", err)
		return _error.HttpErrString("invalid id", response.CM04)
	}

	return nil
}
