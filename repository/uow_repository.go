package repository

import (
	"context"
	"database/sql"
)

//counterfeiter:generate -o ./mocks . UnitOfWorkRepository
type UnitOfWorkRepository interface {
	GetDB() (*sql.DB, error)
	StartTx(ctx context.Context, opts *sql.TxOptions, fn func() error) error
	GetTx() (*sql.Tx, error)
}

func LevelReadCommitted() *sql.TxOptions {
	return &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	}
}
