package repository

import (
	"context"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	CheckUserByEmail(ctx context.Context, email string) (bool, error)
	CheckUserByUsername(ctx context.Context, username string) (bool, error)
	UnitOfWorkRepo
}
