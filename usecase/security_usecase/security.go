package security_usecase

import (
	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
)

type SecurityUsecaseImpl struct {
	userRepo     repository.UserRepository
	securityRepo repository.SecurityRepository
}

func NewSecurityUsecaseImpl(
	userRepo repository.UserRepository,
	securityRepo repository.SecurityRepository,
) usecase.SecurityUsecase {
	return &SecurityUsecaseImpl{
		userRepo:     userRepo,
		securityRepo: securityRepo,
	}
}
