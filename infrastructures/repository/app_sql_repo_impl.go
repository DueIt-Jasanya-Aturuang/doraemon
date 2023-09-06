package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rs/zerolog/log"
	
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
)

type AppRepoSqlImpl struct {
	repository.UnitOfWorkSqlRepo
}

func NewAppRepoSqlImpl(repo repository.UnitOfWorkSqlRepo) repository.AppSqlRepo {
	return &AppRepoSqlImpl{
		UnitOfWorkSqlRepo: repo,
	}
}

func (a *AppRepoSqlImpl) CheckAppByID(ctx context.Context, id string) (bool, error) {
	query := `SELECT id FROM m_app WHERE id = $1`

	conn, err := a.GetConn()
	if err != nil {
		return false, err
	}

	stmt, err := conn.PrepareContext(ctx, query)
	if err != nil {
		log.Err(err).Msg("failed to start prepared context")
		return false, err
	}
	defer func() {
		if errStmt := stmt.Close(); errStmt != nil {
			log.Err(err).Msg("failed to close prepared context")
		}
	}()

	var exists bool
	err = stmt.QueryRowContext(ctx, id).Scan(&exists)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Err(err).Msg("cannot scan query row context")
		}
		return false, err
	}

	if exists {
		return true, nil
	}

	return false, nil
}
