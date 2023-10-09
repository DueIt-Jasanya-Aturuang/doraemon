package user_usecase

import (
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
)

type UserUsecaseImpl struct {
	userRepo repository.UserRepository
	redis    *infra.RedisImpl
}

func NewUserUsecaseImpl(
	userRepo repository.UserRepository,
	redis *infra.RedisImpl,
) usecase.UserUsecase {
	return &UserUsecaseImpl{
		userRepo: userRepo,
		redis:    redis,
	}
}
