package security_usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase/jwt_usecase"
)

func (s *SecurityUsecaseImpl) ReGenerateJWT(ctx context.Context, req *usecase.RequestValidationJWT) (*usecase.ResponseJWT, error) {
	token, err := s.securityRepo.GetByAccessToken(ctx, req.Authorization)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn().Msg("user mencoba untuk request menggunakan token yang lama, kita delete semua tokennya")

			err = s.deletedAllToken(ctx, req.UserID)
		}
		return nil, err
	}

	if token.UserID != req.UserID {
		log.Warn().Msgf("user id jwt and user id header not match | jwt : %s | header : %s", token.UserID, req.UserID)
		return nil, usecase.JwtUserIDAndHeaderUserIDNotMatch
	}

	_, err = jwt_usecase.ClaimsJwtHS256(token.RefreshToken, infra.RefreshTokenKeyHS)
	if err != nil {
		err = s.deletedToken(ctx, token.ID, token.UserID)
		if err != nil {
			return nil, err
		}

		return nil, usecase.InvalidToken
	}

	rtat, err := jwt_usecase.GenerateRTAT(token.UserID, req.AppID, token.RememberMe)
	if err != nil {
		return nil, err
	}

	err = s.securityRepo.StartTx(ctx, repository.LevelReadCommitted(), func() error {
		err = s.securityRepo.Update(ctx, repository.UpdateSecurityParams{
			ID:           token.ID,
			RefreshToken: rtat.RefreshToken,
			AccessToken:  rtat.AcceesToken,
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	resp := &usecase.ResponseJWT{
		Token: rtat.AcceesToken,
	}

	return resp, nil
}
func (s *SecurityUsecaseImpl) GenerateJWT(ctx context.Context, req *usecase.RequestGenerateJWT) (*usecase.ResponseJWT, error) {
	rtat, err := jwt_usecase.GenerateRTAT(req.UserID, req.AppID, req.RememberMe)
	if err != nil {
		return nil, err
	}

	err = s.securityRepo.StartTx(ctx, repository.LevelReadCommitted(), func() error {
		err = s.securityRepo.Create(ctx, rtat)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	resp := &usecase.ResponseJWT{
		Token: rtat.AcceesToken,
	}

	return resp, nil
}
