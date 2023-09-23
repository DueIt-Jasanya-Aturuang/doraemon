package helper

import (
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func BcryptPasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Warn().Msgf("failed to generate password | error : %v", err)
	}

	return string(bytes), err
}

func BcryptPasswordCompare(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
