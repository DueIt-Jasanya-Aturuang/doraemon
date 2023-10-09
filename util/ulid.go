package util

import (
	"github.com/docker/distribution/uuid"
	"github.com/jasanya-tech/jasanya-response-backend-golang/_error"
	"github.com/jasanya-tech/jasanya-response-backend-golang/response"
	"github.com/rs/zerolog/log"
)

func ParseUUID(s string) error {
	if _, err := uuid.Parse(s); err != nil {
		log.Info().Msgf("failed parse ulid | err : %v", err)
		return _error.HttpErrString("invalid id", response.CM04)
	}

	return nil
}
