package otp_usecase

import (
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
)

type OTPUsecaseImpl struct {
	userRepo repository.UserRepository
	redis    *infra.RedisImpl
}

func NewOTPUsecaseImpl(
	userRepo repository.UserRepository,
	redis *infra.RedisImpl,
) usecase.OTPUsecase {
	return &OTPUsecaseImpl{
		userRepo: userRepo,
		redis:    redis,
	}
}
