package repository

import (
	"context"
	"database/sql"
)

type UnitOfWorkSqlRepo interface {
	OpenConn(ctx context.Context) error
	GetConn() (*sql.Conn, error)
	CloseConn()
	StartTx(ctx context.Context, opts *sql.TxOptions) error
	EndTx(err error) error
	GetTx() (*sql.Tx, error)
}
