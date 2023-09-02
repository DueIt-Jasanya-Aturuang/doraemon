package usecase

import (
	"context"
	"time"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/usecase"
)

type SecurityUsecaseImpl struct {
	userRepo repository.UserSqlRepo
	timeout  time.Duration
}

func NewSecurityUsecaseImpl(
	userRepo repository.UserSqlRepo,
	timeout time.Duration,
) usecase.SecurityUsecase {
	return &SecurityUsecaseImpl{
		userRepo: userRepo,
		timeout:  timeout,
	}
}

func (s *SecurityUsecaseImpl) JwtValidateAT(ctx context.Context, req *dto.JwtTokenReq, endpoint string) error {
	// TODO implement me
	panic("implement me")
}

func (s *SecurityUsecaseImpl) JwtGenerateRTAT(ctx context.Context, req *dto.JwtTokenReq) (*dto.JwtTokenResp, error) {
	// TODO implement me
	panic("implement me")
}

func (s *SecurityUsecaseImpl) ResetPassword(ctx context.Context, req *dto.ResetPasswordReq) error {
	// TODO implement me
	panic("implement me")
}

func (s *SecurityUsecaseImpl) ForgottenPassword(ctx context.Context, req *dto.ForgottenPasswordReq) error {
	// TODO implement me
	panic("implement me")
}
