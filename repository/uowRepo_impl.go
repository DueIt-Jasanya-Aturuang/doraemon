package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

type UnitOfWorkRepositoryImpl struct {
	tx *sql.Tx
	db *sql.DB
}

func NewUnitOfWorkRepositoryImpl(db *sql.DB) *UnitOfWorkRepositoryImpl {
	return &UnitOfWorkRepositoryImpl{
		db: db,
	}
}

func (u *UnitOfWorkRepositoryImpl) GetDB() (*sql.DB, error) {
	if u.db == nil {
		err := fmt.Errorf("no Connection Database Available")
		log.Warn().Msg(err.Error())
		return nil, err
	}

	return u.db, nil
}

func (u *UnitOfWorkRepositoryImpl) StartTx(ctx context.Context, opts *sql.TxOptions, fn func() error) error {
	if u.db == nil {
		err := fmt.Errorf("no Connection Database Available")
		log.Warn().Msg(err.Error())
		return err
	}

	tx, err := u.db.BeginTx(ctx, opts)
	if err != nil {
		log.Warn().Msgf(util.LogErrBeginTx, err)
		return err
	}
	context.WithValue(ctx, "tx", "tx")
	u.tx = tx

	err = fn()
	if err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			log.Warn().Msgf(util.LogErrRollback, errRollback)
			return errRollback
		}

		log.Info().Msgf(util.LogInfoRollback, err)
		return err
	}

	if errCommit := tx.Commit(); errCommit != nil {
		log.Warn().Msgf(util.LogErrCommit, errCommit)
		return errCommit
	}

	log.Info().Msgf(util.LogInfoCommit)
	return nil
}

func (u *UnitOfWorkRepositoryImpl) GetTx() (*sql.Tx, error) {
	if u.tx == nil {
		err := fmt.Errorf("no Transaction Available")
		log.Warn().Msg(err.Error())
		return nil, err
	}

	return u.tx, nil
}
