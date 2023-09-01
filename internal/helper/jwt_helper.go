package helper

import (
	"errors"
	"fmt"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"time"
)

func GenerateJwtHS256(jwtModel *model.Jwt) (string, error) {
	timeNow := time.Now()
	timeExp := timeNow.Add(jwtModel.Exp).Unix()

	sub := fmt.Sprintf("%s:%t:%s", jwtModel.UUID, jwtModel.RememberMe, jwtModel.Type)

	tokenParse := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": timeExp,
		"sub": sub,
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

	if claims, ok := tokenParse.Claims.(jwt.MapClaims); ok && tokenParse.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
