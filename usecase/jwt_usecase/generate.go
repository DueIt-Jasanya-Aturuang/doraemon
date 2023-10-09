package jwt_usecase

import (
	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
)

func GenerateRTAT(userID string, appID string, rememberMe bool) (*repository.Security, error) {
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

	resp := &repository.Security{
		UserID:       userID,
		AppID:        appID,
		RememberMe:   rememberMe,
		AcceesToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return resp, nil
}
