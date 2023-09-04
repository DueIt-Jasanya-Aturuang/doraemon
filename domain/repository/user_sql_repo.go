package repository

import (
	"context"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
)

//counterfeiter:generate -o ./../mocks . UserSqlRepo
type UserSqlRepo interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	CheckUserByEmail(ctx context.Context, email string) (bool, error)
	CheckUserByUsername(ctx context.Context, username string) (bool, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByEmailOrUsername(ctx context.Context, emailOrUsername string) (*model.User, error)
	UnitOfWorkSqlRepo
}
