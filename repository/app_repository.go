package repository

import (
	"context"
)

type AppRepository interface {
	CheckAppByID(ctx context.Context, id string) (bool, error)
	UnitOfWorkRepository
}
