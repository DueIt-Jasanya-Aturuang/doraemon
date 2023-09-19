package util

import (
	"encoding/json"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
)

func SerializeProfile(userID string) ([]byte, error) {
	profileReq := dto.ProfileReq{
		UserID: userID,
	}

	profileJson, err := json.Marshal(&profileReq)

	if err != nil {
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
		return nil, err
	}

	return kafkaMsg, nil
}
