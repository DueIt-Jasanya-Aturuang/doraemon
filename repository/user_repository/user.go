package user_repository

import (
	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
)

type UserRepositoryImpl struct {
	repository.UnitOfWorkRepository
}

func NewUserRepositoryImpl(uow repository.UnitOfWorkRepository) repository.UserRepository {
	return &UserRepositoryImpl{
		UnitOfWorkRepository: uow,
	}
}
