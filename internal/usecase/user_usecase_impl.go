package usecase

import (
	"context"
	"time"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/usecase"
)

type UserUsecaseImpl struct {
	userRepo repository.UserSqlRepo
	timeout  time.Duration
}

func NewUserUsecaseImpl(
	userRepo repository.UserSqlRepo,
	timeout time.Duration,
) usecase.UserUsecase {
	return &UserUsecaseImpl{
		userRepo: userRepo,
		timeout:  timeout,
	}
}

func (s *UserUsecaseImpl) ResetPassword(ctx context.Context, req *dto.ResetPasswordReq) error {
	// TODO implement me
	panic("implement me")
}

func (s *UserUsecaseImpl) ForgottenPassword(ctx context.Context, req *dto.ForgottenPasswordReq) error {
	// TODO implement me
	panic("implement me")
}
