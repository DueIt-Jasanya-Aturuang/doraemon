package repository

import (
	"context"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate -o ./../mocks . AccessSqlRepo
type AccessSqlRepo interface {
	CreateAccess(ctx context.Context, access *model.Access) (*model.Access, error)
	GetAccessByUserIDAndAppID(ctx context.Context, userID string, appID string) (*model.Access, error)
	UnitOfWorkSqlRepo
}
