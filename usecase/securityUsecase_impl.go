package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase/helper"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

type SecurityUsecaseImpl struct {
	userRepo     domain.UserRepository
	securityRepo domain.SecurityRepository
}

func NewSecurityUsecaseImpl(
	userRepo domain.UserRepository,
	securityRepo domain.SecurityRepository,
) domain.SecurityUsecase {
	return &SecurityUsecaseImpl{
		userRepo:     userRepo,
		securityRepo: securityRepo,
	}
}

func (s *SecurityUsecaseImpl) JwtValidation(ctx context.Context, req *domain.RequestJwtToken) (bool, error) {
	claims, err := helper.ClaimsJwtHS256(req.Authorization, infra.AccessTokenKeyHS)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return true, nil
		}
		log.Warn().Msgf(util.LogErrFailedClaimJwt, err, req.Authorization)
		return false, InvalidToken
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		log.Warn().Msgf("cannot assetion claims sub to string | data : %v", claims["sub"])
		return false, InvalidToken
	}

	if userID != req.UserId {
		log.Warn().Msgf("user id jwt and user id header not match | jwt : %s | header : %s", userID, req.UserId)
		return false, JwtUserIDAndHeaderUserIDNotMatch
	}

	if !req.ActivasiHeader {
		activasi, err := s.userRepo.CheckActivasiUser(ctx, req.UserId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return false, InvalidUserID
			}
			return false, err
		}

		if !activasi {
			return false, UserIsNotActivited
		}

	}

	token, err := s.securityRepo.GetByAccessToken(ctx, req.Authorization)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn().Msg("user mencoba untuk request menggunakan token yang lama, kita delete semua tokennya")

			err = s.deletedAllToken(ctx, req.UserId)
		}
		return false, err
	}

	if req.AppId != token.AppID {
		return false, JwtAppIDAndHeaderAppIDNotMatch
	}

	return false, nil
}

func (s *SecurityUsecaseImpl) JwtGenerate(ctx context.Context, req *domain.RequestJwtToken) (*domain.ResponseJwtToken, error) {
	token, err := s.securityRepo.GetByAccessToken(ctx, req.Authorization)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn().Msg("user mencoba untuk request menggunakan token yang lama, kita delete semua tokennya")

			err = s.deletedAllToken(ctx, req.UserId)
		}
		return nil, err
	}

	if token.UserID != req.UserId {
		log.Warn().Msgf("user id jwt and user id header not match | jwt : %s | header : %s", token.UserID, req.UserId)
		return nil, JwtUserIDAndHeaderUserIDNotMatch
	}

	_, err = helper.ClaimsJwtHS256(token.RefreshToken, infra.RefreshTokenKeyHS)
	if err != nil {
		err = s.deletedToken(ctx, token.ID, token.UserID)
		if err != nil {
			return nil, err
		}

		return nil, InvalidToken
	}

	rtat, err := helper.GenerateRTAT(token.UserID, req.AppId, token.RememberMe)
	if err != nil {
		return nil, err
	}

	err = s.securityRepo.StartTx(ctx, helper.LevelReadCommitted(), func() error {
		err = s.securityRepo.Update(ctx, token.ID, rtat.RefreshToken, rtat.AcceesToken)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	resp := &domain.ResponseJwtToken{
		Token: rtat.AcceesToken,
	}

	return resp, nil
}

func (s *SecurityUsecaseImpl) Logout(ctx context.Context, req *domain.RequestLogout) error {
	token, err := s.securityRepo.GetByAccessToken(ctx, req.Token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn().Msg("user mencoba untuk logout menggunakan token yang lama")
			return nil
		}
		return err
	}

	err = s.deletedToken(ctx, token.ID, token.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (s *SecurityUsecaseImpl) deletedAllToken(ctx context.Context, userID string) error {
	err := s.securityRepo.StartTx(ctx, helper.LevelReadCommitted(), func() error {
		err := s.securityRepo.DeleteAllByUserID(ctx, userID)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return InvalidToken
}

func (s *SecurityUsecaseImpl) deletedToken(ctx context.Context, id int, userID string) error {
	err := s.securityRepo.StartTx(ctx, helper.LevelReadCommitted(), func() error {
		err := s.securityRepo.Delete(ctx, id, userID)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
