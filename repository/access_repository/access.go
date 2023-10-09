package access_repository

import (
	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
)

type AccessRepositoryImpl struct {
	repository.UnitOfWorkRepository
}

func NewAccessRepositoryImpl(uow repository.UnitOfWorkRepository) repository.AccessRepository {
	return &AccessRepositoryImpl{
		UnitOfWorkRepository: uow,
	}
}
