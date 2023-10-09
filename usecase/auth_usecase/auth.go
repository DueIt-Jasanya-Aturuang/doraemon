package auth_usecase

import (
	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
)

type AuthUsecaseImpl struct {
	userRepo          repository.UserRepository
	accessRepo        repository.AccessRepository
	apiServiceUsecase usecase.ApiServiceUsecase
	securityUsecase   usecase.SecurityUsecase
}

func NewAuthUsecaseImpl(
	userRepo repository.UserRepository,
	accessRepo repository.AccessRepository,
	apiServiceUsecase usecase.ApiServiceUsecase,
	securityUsecase usecase.SecurityUsecase,
) usecase.AuthUsecase {
	return &AuthUsecaseImpl{
		userRepo:          userRepo,
		accessRepo:        accessRepo,
		apiServiceUsecase: apiServiceUsecase,
		securityUsecase:   securityUsecase,
	}

}
