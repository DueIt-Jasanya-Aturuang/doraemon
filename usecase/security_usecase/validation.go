package security_usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase/jwt_usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

func (s *SecurityUsecaseImpl) ValidationJWT(ctx context.Context, req *usecase.RequestValidationJWT) (bool, error) {
	claims, err := jwt_usecase.ClaimsJwtHS256(req.Authorization, infra.AccessTokenKeyHS)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return true, nil
		}
		log.Warn().Msgf(util.LogErrFailedClaimJwt, err, req.Authorization)
		return false, usecase.InvalidToken
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		log.Warn().Msgf("cannot assetion claims sub to string | data : %v", claims["sub"])
		return false, usecase.InvalidToken
	}

	if userID != req.UserID {
		log.Warn().Msgf("user_repository id jwt_usecase and user_repository id header not match | jwt_usecase : %s | header : %s", userID, req.UserID)
		return false, usecase.JwtUserIDAndHeaderUserIDNotMatch
	}

	if !req.ActivasiHeader {
		activasi, err := s.userRepo.CheckActivasiUser(ctx, req.UserID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return false, usecase.InvalidUserID
			}
			return false, err
		}

		if !activasi {
			return false, usecase.UserIsNotActivited
		}

	}

	token, err := s.securityRepo.GetByAccessToken(ctx, req.Authorization)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn().Msg("user_repository mencoba untuk request menggunakan token yang lama, kita delete semua tokennya")

			err = s.deletedAllToken(ctx, req.UserID)
		}
		return false, err
	}

	if req.AppID != token.AppID {
		return false, usecase.JwtAppIDAndHeaderAppIDNotMatch
	}

	return false, nil
}
