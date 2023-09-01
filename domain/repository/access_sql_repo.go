package repository

import (
	"context"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
)

type AccessSqlRepo interface {
	CreateAccess(ctx context.Context, access *model.Access) (*model.Access, error)
	GetAccessByUserIDAndAppID(ctx context.Context, userID string, appID string) (*model.Access, error)
	UnitOfWorkSqlRepo
}
