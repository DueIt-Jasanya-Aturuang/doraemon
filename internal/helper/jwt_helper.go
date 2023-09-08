package helper

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
)

func GenerateJwtHS256(jwtModel *model.Jwt) (string, error) {
	timeNow := time.Now()
	timeExp := timeNow.Add(jwtModel.Exp).Unix()

	tokenParse := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": timeExp,
		"sub": jwtModel.UserID,
	})

	tokenStr, err := tokenParse.SignedString([]byte(jwtModel.Key))
	if err != nil {
		log.Err(err).Msg("failed signing token string hs 256")
		return "", err
	}

	return tokenStr, err
}

func ClaimsJwtHS256(tokenStr, key string) (map[string]any, error) {
	tokenParse, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Warn().Msgf("unexpected signing method : %v", t.Header["alg"])
			return nil, errors.New("invalid token")
		}
		return []byte(key), nil
	})

	if err != nil {
		return nil, err
	}

	claims, _ := tokenParse.Claims.(jwt.MapClaims)

	return claims, nil
}
