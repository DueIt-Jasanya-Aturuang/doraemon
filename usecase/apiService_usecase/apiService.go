package apiService_usecase

import (
	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
)

type ApiServiceUsecaseImpl struct {
	apiServiceRepo repository.ApiServiceRepository
}

func NewApiServiceUsecaseImpl(
	apiServiceRepo repository.ApiServiceRepository,
) usecase.ApiServiceUsecase {
	return &ApiServiceUsecaseImpl{
		apiServiceRepo: apiServiceRepo,
	}
}
