package usecase

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/satori/go.uuid"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infrastructures/config"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/helper"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/encryption"
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

func (s *SecurityUsecaseImpl) JwtValidateAT(ctx context.Context, req *dto.JwtTokenReq, endpoint string) (string, bool, error) {
	claims, err := helper.ClaimsJwtHS256(req.Authorization, config.AccessTokenKeyHS)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", true, nil
		}
		return "", false, _error.ErrString("INVALID YOUR TOKEN", http.StatusUnauthorized)
	}

	sub, err := encryption.DecryptStringCFB(claims["sub"].(string), config.AesCFB)
	if err != nil {
		return "", false, _error.ErrString("INVALID YOUR TOKEN", http.StatusUnauthorized)
	}

	subArray := strings.Split(sub, ":")
	if len(subArray) != 3 {
		return "", false, _error.ErrString("INVALID YOUR TOKEN", http.StatusUnauthorized)
	}

	if subArray[2] != "access-token" {
		return "", false, _error.ErrString("INVALID YOUR TOKEN", http.StatusUnauthorized)
	}

	if subArray[1] != req.UserId {
		return "", false, _error.ErrString("INVALID YOUR TOKEN", http.StatusUnauthorized)
	}

	err = s.securityRepo.OpenConn(ctx)
	if err != nil {
		return "", false, _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer s.securityRepo.CloseConn()

	getRT, err := s.securityRepo.GetTokenByIDAndUserID(ctx, subArray[0], req.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err := s.securityRepo.DeleteAllTokenByUserID(ctx, req.UserId)
			if err != nil {
				return "", false, _error.ErrStringDefault(http.StatusInternalServerError)
			}
			return "", false, _error.ErrString("INVALID YOUR TOKEN", http.StatusUnauthorized)
		}
		return "", false, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	if req.AppId != getRT.AppID {
		return "", false, _error.ErrStringDefault(http.StatusForbidden)
	}

	if strings.Contains(endpoint, "/activasi-account") {
		activasi, err := s.userRepo.CheckActivasiUserByID(ctx, req.UserId)
		if err != nil {
			return "", false, _error.ErrStringDefault(http.StatusInternalServerError)
		}
		if !activasi {
			return "", false, _error.ErrString("akun anda tidak aktif, silahkan aktifkan akun anda", http.StatusForbidden)
		}
	}

	return getRT.Token, false, nil
}

func (s *SecurityUsecaseImpl) JwtGenerateRTAT(ctx context.Context, req *dto.JwtTokenReq) (tokenResp *dto.JwtTokenResp, err error) {
	claims, err := helper.ClaimsJwtHS256(req.Authorization, config.AccessTokenKeyHS)
	if err != nil {
		return nil, _error.ErrString("INVALID YOUR TOKEN", http.StatusUnauthorized)
	}

	sub, err := encryption.DecryptStringCFB(claims["sub"].(string), config.AesCFB)
	if err != nil {
		return nil, _error.ErrString("INVALID YOUR TOKEN", http.StatusUnauthorized)
	}

	subArray := strings.Split(sub, ":")
	tokenID := subArray[0]

	err = s.securityRepo.OpenConn(ctx)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer s.securityRepo.CloseConn()

	getRT, err := s.securityRepo.GetTokenByIDAndUserID(ctx, tokenID, req.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err := s.securityRepo.DeleteAllTokenByUserID(ctx, req.UserId)
			if err != nil {
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
	defer func() {
		errEndTx := s.userRepo.EndTx(err)
		if errEndTx != nil {
			err = _error.ErrStringDefault(http.StatusInternalServerError)
		}
	}()

	err = s.securityRepo.DeleteToken(ctx, tokenID, req.UserId)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	newTokenID := uuid.NewV4().String()
	var jwtModel *model.Jwt

	jwtModelAT := jwtModel.AccessTokenDefault(newTokenID, req.UserId, getRT.RememberMe)
	accessToken, err := helper.GenerateJwtHS256(jwtModelAT)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	jwtModelRT := jwtModel.RefreshTokenDefault(newTokenID, req.UserId, getRT.RememberMe)
	refreshToken, err := helper.GenerateJwtHS256(jwtModelRT)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	tokenUpdate := &model.TokenUpdate{
		ID:    newTokenID,
		OldID: tokenID,
		Token: refreshToken,
	}
	err = s.securityRepo.UpdateToken(ctx, tokenUpdate)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	tokenResp = &dto.JwtTokenResp{
		RememberMe: getRT.RememberMe,
		Token:      accessToken,
		Exp:        config.AccessTokenKeyExpHS,
	}

	return tokenResp, nil
}

func (s *SecurityUsecaseImpl) JwtRegistredRTAT(ctx context.Context, req *dto.JwtRegisteredTokenReq) (tokenResp *dto.JwtTokenResp, err error) {
	err = s.userRepo.OpenConn(ctx)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer s.userRepo.CloseConn()

	tokenID := uuid.NewV4().String()
	var jwtModel *model.Jwt

	jwtModelAT := jwtModel.AccessTokenDefault(tokenID, req.UserId, req.RememberMe)
	accessToken, err := helper.GenerateJwtHS256(jwtModelAT)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	jwtModelRT := jwtModel.RefreshTokenDefault(tokenID, req.UserId, req.RememberMe)
	refreshToken, err := helper.GenerateJwtHS256(jwtModelRT)
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

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
		ID:         tokenID,
		UserID:     req.UserId,
		AppID:      req.AppId,
		Token:      refreshToken,
		RememberMe: req.RememberMe,
	})
	if err != nil {
		return nil, _error.ErrStringDefault(http.StatusInternalServerError)
	}

	tokenResp = &dto.JwtTokenResp{
		RememberMe: req.RememberMe,
		Token:      accessToken,
		Exp:        jwtModelAT.Exp,
	}

	return tokenResp, nil
}
