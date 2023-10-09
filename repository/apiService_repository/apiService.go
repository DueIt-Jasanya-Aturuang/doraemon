package apiService_repository

import (
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
)

type ApiServiceRepositoryImpl struct {
	baseUrlDueit string
}

func NewApiServiceRepositoryImpl() repository.ApiServiceRepository {
	return &ApiServiceRepositoryImpl{
		baseUrlDueit: infra.BaseUrlDueitAccountService,
	}
}
