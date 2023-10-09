package security_repository

import (
	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
)

type SecurityRepositoryImpl struct {
	repository.UnitOfWorkRepository
}

func NewSecurityRepositoryImpl(
	uow repository.UnitOfWorkRepository,
) repository.SecurityRepository {
	return &SecurityRepositoryImpl{
		UnitOfWorkRepository: uow,
	}
}
