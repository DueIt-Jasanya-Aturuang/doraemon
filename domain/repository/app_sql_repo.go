package repository

import "context"

//counterfeiter:generate -o ./../mocks . AppSqlRepo
type AppSqlRepo interface {
	CheckAppByID(ctx context.Context, id string) (bool, error)
	UnitOfWorkSqlRepo
}
