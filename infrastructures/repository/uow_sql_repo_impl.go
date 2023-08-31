package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
	"github.com/rs/zerolog/log"
)

type UnitOfWorkImpl struct {
	db   *sql.DB
	tx   *sql.Tx
	conn *sql.Conn
}

func NewUnitOfWorkImpl(db *sql.DB) repository.UnitOfWorkRepo {
	return &UnitOfWorkImpl{
		db: db,
	}
}

func (u *UnitOfWorkImpl) OpenConn(ctx context.Context) error {
	conn, err := u.db.Conn(ctx)
	if err != nil {
		log.Err(err).Msg("Cannot Open Database")
		return err
	}

	u.conn = conn

	return nil
}

func (u *UnitOfWorkImpl) GetConn() (*sql.Conn, error) {
	if u.conn == nil {
		err := fmt.Errorf("no connection available")
		log.Err(err).Msg("no connection available")
		return nil, err
	}

	return u.conn, nil
}

func (u *UnitOfWorkImpl) CloseConn() {
	err := u.conn.Close()
	if err != nil {
		log.Err(err).Msg("Cannot close Database")
	}
}

func (u *UnitOfWorkImpl) StartTx(ctx context.Context, opts *sql.TxOptions) error {
	if u.conn == nil {
		err := fmt.Errorf("no connection available")
		log.Err(err).Msg(err.Error())
		return err
	}
	tx, err := u.conn.BeginTx(ctx, opts)
	if err != nil {
		log.Err(err).Msg("cannot start transaction")
		return err
	}

	u.tx = tx

	return nil
}

func (u *UnitOfWorkImpl) EndTx(err error) error {
	if err != nil && !errors.Is(err, sql.ErrTxDone) {
		if errRollback := u.tx.Rollback(); errRollback != nil {
			log.Err(errRollback).Msgf("cannot rollback transaction from error : %v", err)
			return errRollback
		}
	} else {
		if errCommit := u.tx.Commit(); errCommit != nil {
			log.Err(errCommit).Msg("cannot commit transaction")
			return errCommit
		}
	}

	return nil
}

func (u *UnitOfWorkImpl) GetTx() (*sql.Tx, error) {
	if u.tx == nil {
		err := fmt.Errorf("no transaction available")
		log.Err(err).Msg("no transaction available")
		return nil, err
	}

	return u.tx, nil
}
