package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
	_msg "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/msg"
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
	query := `SELECT EXISTS(SELECT 1 FROM m_apps WHERE id = $1 AND deleted_at IS NULL)`
	conn, err := a.GetConn()
	if err != nil {
		return false, err
	}

	stmt, err := conn.PrepareContext(ctx, query)
	if err != nil {
		log.Err(err).Msg(_msg.LogErrStartPrepareContext)
		return false, err
	}
	defer func() {
		if errStmt := stmt.Close(); errStmt != nil {
			log.Err(err).Msg(_msg.LogErrClosePrepareContext)
		}
	}()

	var exists bool
	err = stmt.QueryRowContext(ctx, id).Scan(&exists)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Err(err).Msg(_msg.LogErrQueryRowContextScan)
		}
		return false, err
	}

	if exists {
		return true, nil
	}

	return false, nil
}
