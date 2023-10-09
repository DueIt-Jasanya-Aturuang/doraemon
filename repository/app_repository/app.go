package app_repository

import (
	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
)

type AppRepositoryImpl struct {
	repository.UnitOfWorkRepository
}

func NewAppRepositoryImpl(uow repository.UnitOfWorkRepository) repository.AppRepository {
	return &AppRepositoryImpl{
		UnitOfWorkRepository: uow,
	}
}
