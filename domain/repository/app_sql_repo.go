package repository

import "context"

type AppRepo interface {
	CheckAppByID(ctx context.Context, id string) (bool, error)
	UnitOfWorkRepo
}
