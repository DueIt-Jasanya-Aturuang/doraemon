package jwt_usecase

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

func ClaimsJwtHS256(tokenStr, key string) (map[string]any, error) {
	tokenParse, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Warn().Msgf("unexpected signing method : %v", t.Header["alg"])
			return nil, errors.New("invalid Token")
		}
		return []byte(key), nil
	})

	if err != nil {
		return nil, err
	}

	claims, _ := tokenParse.Claims.(jwt.MapClaims)

	return claims, nil
}
