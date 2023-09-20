package util

import (
	"crypto/rand"
	"math/big"

	"github.com/rs/zerolog/log"
)

func charseByte(length int, charset string) ([]byte, error) {
	b := make([]byte, length)
	maxCharset := big.NewInt(int64(len(charset)))

	for i := range b {
		n, err := rand.Int(rand.Reader, maxCharset)
		if err != nil {
			log.Err(err).Msg("CANNOT GENERATE RAND INT")
			return nil, err
		}
		b[i] = charset[n.Int64()]
	}

	return b, nil
}

func RandomChar(length int) (string, error) {
	randomByte, err := charseByte(length, charset)
	if err != nil {
		log.Err(err).Msg("failed generate random character")
		return "", err
	}

	return string(randomByte), err
}
