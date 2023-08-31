package repository

import (
	"context"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
	"github.com/rs/zerolog/log"
)

type AppRepoImpl struct {
	repository.UnitOfWorkRepo
}

func NewAppRepoImpl(repo repository.UnitOfWorkRepo) repository.AppRepo {
	return &AppRepoImpl{
		UnitOfWorkRepo: repo,
	}
}

func (a *AppRepoImpl) CheckAppByID(ctx context.Context, id string) (bool, error) {
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

	var exists bool
	err = stmt.QueryRowContext(ctx, id).Scan(&exists)
	if err != nil {
		log.Err(err).Msg("failed to query row context prepared statement")
		return false, err
	}

	if exists {
		return true, nil
	}

	return false, nil
}
