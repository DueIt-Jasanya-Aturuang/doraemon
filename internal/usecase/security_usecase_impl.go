package usecase

import (
	"context"
	"database/sql"
	"net/http"

	uuid "github.com/satori/go.uuid"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/usecase"
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

func (s *SecurityUsecaseImpl) JwtValidateAT(ctx context.Context, req *dto.JwtTokenReq, endpoint string) (*dto.JwtTokenResp, error) {
	// TODO implement me
	panic("implement me")
}

func (s *SecurityUsecaseImpl) JwtGenerateRTAT(ctx context.Context, req *dto.JwtTokenReq) (*dto.JwtTokenResp, error) {
	// TODO implement me
	panic("implement me")
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
		ID:     tokenID,
		UserID: req.UserId,
		AppID:  req.AppId,
		Token:  refreshToken,
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
