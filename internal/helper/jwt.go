package helper

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra"
)

func GenerateJwtHS256(jwtModel *Jwt) (string, error) {
	timeNow := time.Now()
	timeExp := timeNow.Add(jwtModel.Exp).Unix()

	tokenParse := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": timeExp,
		"sub": jwtModel.UserID,
	})

	tokenStr, err := tokenParse.SignedString([]byte(jwtModel.Key))
	if err != nil {
		log.Warn().Msgf("failed signing token string hs 256 | err : %v", err)
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

func GenerateRTAT(userID string, appID string, rememberMe bool) (*domain.Token, error) {
	var jwtModel *Jwt

	jwtModelAT := jwtModel.AccessTokenDefault(userID)
	accessToken, err := GenerateJwtHS256(jwtModelAT)
	if err != nil {
		log.Warn().Msgf("failed generate at | err : %v", err)
		return nil, err
	}

	jwtModelRT := jwtModel.RefreshTokenDefault(userID, rememberMe)
	refreshToken, err := GenerateJwtHS256(jwtModelRT)
	if err != nil {
		log.Warn().Msgf("failed generate rt | err : %v", err)
		return nil, err
	}

	resp := &domain.Token{
		UserID:       userID,
		AppID:        appID,
		RememberMe:   rememberMe,
		AcceesToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return resp, nil
}

type Jwt struct {
	UserID string
	Key    string
	Exp    time.Duration
}

func (j *Jwt) AccessTokenDefault(userID string) *Jwt {
	return &Jwt{
		UserID: userID,
		Key:    infra.AccessTokenKeyHS,
		Exp:    infra.AccessTokenKeyExpHS,
	}
}

func (j *Jwt) RefreshTokenDefault(userID string, rememberMe bool) *Jwt {
	var exp time.Duration
	if rememberMe {
		exp = infra.RememberMeTokenExp
	} else {
		exp = infra.RefreshTokenKeyExpHS
	}

	return &Jwt{
		UserID: userID,
		Key:    infra.RefreshTokenKeyHS,
		Exp:    exp,
	}
}

func (j *Jwt) ForgotPasswordTokenDefault(userID string) *Jwt {
	return &Jwt{
		UserID: userID,
		Key:    infra.DefaultKey,
		Exp:    infra.ForgotPasswordTokenExp,
	}
}
