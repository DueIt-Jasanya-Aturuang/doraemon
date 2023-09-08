package usecase

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infrastructures/config"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/helper"
	_error "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/error"
)

type SecurityUsecaseImpl struct {
	userRepo     repository.UserSqlRepo
	securityRepo repository.SecuritySqlRepo
}

func NewSecurityUsecaseImpl(
	userRepo repository.UserSqlRepo,
	securityRepo repository.SecuritySqlRepo,
) usecase.SecurityUsecase {
	return &SecurityUsecaseImpl{
		userRepo:     userRepo,
		securityRepo: securityRepo,
	}
}

func (s *SecurityUsecaseImpl) JwtValidateAT(ctx context.Context, req *dto.JwtTokenReq, endpoint string) (bool, error) {
	claims, err := helper.ClaimsJwtHS256(req.Authorization, config.AccessTokenKeyHS)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return true, nil
		}
		log.Err(err).Msg("failed claim jwt token")
		return false, _error.ErrString("INVALID YOUR TOKEN", http.StatusUnauthorized)
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		log.Warn().Msgf("failed to assertion sub jwt ke string")
		return false, _error.ErrString("INVALID YOUR TOKEN", http.StatusUnauthorized)
	}

	if req.UserId != userID {
		log.Warn().Msg("request user id dan header user id tidak sama")
		return false, _error.ErrString("INVALID YOUR TOKEN", http.StatusUnauthorized)
	}

	err = s.securityRepo.OpenConn(ctx)
	if err != nil {
		return false, _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer s.securityRepo.CloseConn()

	getToken, err := s.securityRepo.GetTokenByAT(ctx, req.Authorization)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn().Msg("user mencoba untuk request menggunakan token yang lama, kita delete semua tokennya")
			if err := s.securityRepo.DeleteAllTokenByUserID(ctx, req.UserId); err != nil {
				return false, _error.ErrStringDefault(http.StatusInternalServerError)
			}
			return false, _error.ErrString("INVALID YOUR TOKEN", http.StatusUnauthorized)
		}
		return false, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	if req.AppId != getToken.AppID {
		log.Warn().Msg("app id in header and app id di token database tidak sama")
		return false, _error.ErrStringDefault(http.StatusForbidden)
	}

	if !strings.Contains(endpoint, "/activasi-account") {
		activasi, err := s.userRepo.CheckActivasiUserByID(ctx, req.UserId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return false, _error.ErrStringDefault(http.StatusNotFound)
			}
			return false, _error.ErrStringDefault(http.StatusInternalServerError)
		}
		if !activasi {
			log.Warn().Msg("user belum melakukan activasi tetapi mencoba request ke endpoint lain")
			return false, _error.ErrString("akun anda tidak aktif, silahkan aktifkan akun anda", http.StatusForbidden)
		}
	}

	return false, nil
}

func (s *SecurityUsecaseImpl) JwtGenerateRTAT(ctx context.Context, req *dto.JwtTokenReq) (tokenResp *dto.JwtTokenResp, err error) {
	err = s.securityRepo.OpenConn(ctx)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer s.securityRepo.CloseConn()

	getToken, err := s.securityRepo.GetTokenByAT(ctx, req.Authorization)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn().Msg("user mencoba untuk request menggunakan token yang lama, kita delete semua tokennya")
			if err := s.securityRepo.DeleteAllTokenByUserID(ctx, req.UserId); err != nil {
				return nil, _error.ErrStringDefault(http.StatusInternalServerError)
			}
			return nil, _error.ErrString("INVALID YOUR TOKEN", http.StatusUnauthorized)
		}
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	err = s.securityRepo.StartTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	err = s.securityRepo.DeleteToken(ctx, getToken.ID, req.UserId)
	if err != nil {
		errEndTx := s.userRepo.EndTx(err)
		if errEndTx != nil {
			err = _error.ErrStringDefault(http.StatusInternalServerError)
		}
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	_, err = helper.ClaimsJwtHS256(getToken.RefreshToken, config.RefreshTokenKeyHS)
	if err != nil {
		errEndTx := s.userRepo.EndTx(nil)
		if errEndTx != nil {
			err = _error.ErrStringDefault(http.StatusInternalServerError)
		}
		if errors.Is(err, jwt.ErrTokenExpired) {
			log.Warn().Msg("failed claim refresh token user")
		} else {
			log.Err(err).Msg("failed claim refresh token user")
		}
		return nil, _error.ErrString("INVALID YOUR TOKEN", http.StatusUnauthorized)
	}

	var jwtModel *model.Jwt

	jwtModelAT := jwtModel.AccessTokenDefault(req.UserId)
	accessToken, err := helper.GenerateJwtHS256(jwtModelAT)
	if err != nil {
		errEndTx := s.userRepo.EndTx(err)
		if errEndTx != nil {
			err = _error.ErrStringDefault(http.StatusInternalServerError)
		}
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	jwtModelRT := jwtModel.RefreshTokenDefault(req.UserId, getToken.RememberMe)
	refreshToken, err := helper.GenerateJwtHS256(jwtModelRT)
	if err != nil {
		errEndTx := s.userRepo.EndTx(err)
		if errEndTx != nil {
			err = _error.ErrStringDefault(http.StatusInternalServerError)
		}
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	err = s.securityRepo.UpdateToken(ctx, getToken.ID, refreshToken, accessToken)
	if err != nil {
		errEndTx := s.userRepo.EndTx(err)
		if errEndTx != nil {
			err = _error.ErrStringDefault(http.StatusInternalServerError)
		}
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	errEndTx := s.userRepo.EndTx(err)
	if errEndTx != nil {
		err = _error.ErrStringDefault(http.StatusInternalServerError)
	}

	tokenResp = &dto.JwtTokenResp{
		Token: accessToken,
	}

	return tokenResp, nil
}

func (s *SecurityUsecaseImpl) JwtRegistredRTAT(ctx context.Context, req *dto.JwtRegisteredTokenReq) (tokenResp *dto.JwtTokenResp, err error) {
	var jwtModel *model.Jwt

	jwtModelAT := jwtModel.AccessTokenDefault(req.UserId)
	accessToken, err := helper.GenerateJwtHS256(jwtModelAT)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	jwtModelRT := jwtModel.RefreshTokenDefault(req.UserId, req.RememberMe)
	refreshToken, err := helper.GenerateJwtHS256(jwtModelRT)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	err = s.userRepo.OpenConn(ctx)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer s.userRepo.CloseConn()

	err = s.userRepo.StartTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer func() {
		errEndTx := s.userRepo.EndTx(err)
		if errEndTx != nil {
			err = _error.ErrStringDefault(http.StatusInternalServerError)
			tokenResp = nil
		}
	}()

	err = s.securityRepo.CreateToken(ctx, &model.Token{
		UserID:       req.UserId,
		AppID:        req.AppId,
		RefreshToken: refreshToken,
		RememberMe:   req.RememberMe,
		AcceesToken:  accessToken,
	})
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	tokenResp = &dto.JwtTokenResp{
		Token: accessToken,
	}

	return tokenResp, nil
}

func (s *SecurityUsecaseImpl) Logout(ctx context.Context, req *dto.LogoutReq) (err error) {
	err = s.securityRepo.OpenConn(ctx)
	if err != nil {
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer s.securityRepo.CloseConn()

	getToken, err := s.securityRepo.GetTokenByAT(ctx, req.Token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn().Msg("user mencoba untuk logout menggunakan token yang lama")
			return nil
		}
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}

	err = s.securityRepo.StartTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer func() {
		errEndTx := s.securityRepo.EndTx(err)
		if errEndTx != nil {
			err = _error.ErrStringDefault(http.StatusInternalServerError)
		}
	}()
	err = s.securityRepo.DeleteToken(ctx, getToken.ID, req.UserID)
	if err != nil {
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}

	return nil
}
