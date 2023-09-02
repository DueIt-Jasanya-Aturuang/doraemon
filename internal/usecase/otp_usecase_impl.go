package usecase

import (
	"context"
	"time"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/usecase"
)

type OTPUsecaseImpl struct {
	userRepo repository.UserSqlRepo
	timeout  time.Duration
}

func NewOTPUsecaseImpl(
	userRepo repository.UserSqlRepo,
	timeout time.Duration,
) usecase.OTPUsecase {
	return &OTPUsecaseImpl{
		userRepo: userRepo,
		timeout:  timeout,
	}
}

func (o *OTPUsecaseImpl) OTPGenerate(ctx context.Context, req *dto.OTPGenerateReq) error {
	// TODO implement me
	panic("implement me")
}

func (o *OTPUsecaseImpl) OTPValidation(ctx context.Context, req *dto.OTPValidationReq) error {
	// TODO implement me
	panic("implement me")
}
