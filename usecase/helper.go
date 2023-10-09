package usecase

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

func SerializeMsgKafka(otp, email, typeReq string) ([]byte, error) {
	msg := map[string]string{
		"value": otp,
		"to":    email,
		"type":  typeReq,
	}

	kafkaMsg, err := json.Marshal(msg)

	if err != nil {
		log.Warn().Msgf(util.LogErrMarshal, msg, err)
		return nil, err
	}

	return kafkaMsg, nil
}

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

func FormatEmail(email string) string {
	emailArr := strings.Split(email, "@")
	if len(emailArr) != 2 {
		log.Warn().Msgf("email User tidak valid | total : %d", len(emailArr))
		return email
	}
	emailString := fmt.Sprintf("%c••••%c@%s", emailArr[0][0], emailArr[0][len(emailArr[0])-1], emailArr[1])
	return emailString
}
