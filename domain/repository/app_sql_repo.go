package repository

import "context"

type AppSqlRepo interface {
	CheckAppByID(ctx context.Context, id string) (bool, error)
	UnitOfWorkSqlRepo
}
