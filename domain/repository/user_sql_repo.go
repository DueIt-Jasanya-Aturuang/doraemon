package repository

import (
	"context"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	CheckEmailUser(ctx context.Context, email string) (bool, error)
	CheckUsernameUser(ctx context.Context, username string) (bool, error)
	UnitOfWorkRepo
}
