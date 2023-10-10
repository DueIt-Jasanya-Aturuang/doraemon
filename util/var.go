package util

import (
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
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

var NewUUID = uuid.NewV4().String()

func ParseUUID(u string) error {
	if _, err := ulid.Parse(u); err != nil {
		log.Info().Msgf("failed parse uuid | err : %v", err)
		return err
	}

	return nil
}
