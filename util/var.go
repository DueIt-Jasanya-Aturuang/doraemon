package util

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	uuidSatori "github.com/satori/go.uuid"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const ActivasiAccount = "activasi-account"
const ForgotPasswordLink = "forgot-password-link"
const ForgotPassword = "forgot-password"
const AppIDHeader = "App-ID"
const TypeHeader = "Type"
const UserIDHeader = "User-ID"
const AuthorizationHeader = "Authorization"
const ActivasiHeader = "Activasi"
const DeviceTypeWeb = "web"
const DeviceTypeMobile = "mobile"

var Endpoint = []string{
	"/auth/login",
	"/auth/register",
}

var NewUUID = uuidSatori.NewV4().String()

func ParseUUID(u string) error {
	if _, err := uuid.Parse(u); err != nil {
		log.Info().Msgf("failed parse uuid | err : %v", err)
		return err
	}

	return nil
}
