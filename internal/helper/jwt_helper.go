package helper

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
	_error "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/error"
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

func GenerateRTAT(userID string, appID string, rememberMe bool) (*model.Token, error) {
	var jwtModel *model.Jwt

	jwtModelAT := jwtModel.AccessTokenDefault(userID)
	accessToken, err := GenerateJwtHS256(jwtModelAT)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	jwtModelRT := jwtModel.RefreshTokenDefault(userID, rememberMe)
	refreshToken, err := GenerateJwtHS256(jwtModelRT)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	resp := &model.Token{
		UserID:       userID,
		AppID:        appID,
		RememberMe:   rememberMe,
		AcceesToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return resp, nil
}
