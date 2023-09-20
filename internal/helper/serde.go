package helper

import (
	"encoding/json"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

func EndPointMarshal() string {
	endpoint := []string{
		"/auth/login",
		"/auth/register",
	}

	marshal, err := json.Marshal(endpoint)
	if err != nil {
		log.Warn().Msgf(util.LogErrMarshal, endpoint, err)
		return ""
	}

	return string(marshal)
}

func SerializeProfile(userID string) ([]byte, error) {
	profileReq := domain.RequestCreateProfile{
		UserID: userID,
	}

	profileJson, err := json.Marshal(&profileReq)

	if err != nil {
		log.Warn().Msgf(util.LogErrMarshal, userID, err)
		return nil, err
	}

	return profileJson, err
}

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
