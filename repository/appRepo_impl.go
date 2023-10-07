package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

type AppRepositoryImpl struct {
	domain.UnitOfWorkRepository
}

func NewAppRepositoryImpl(repo domain.UnitOfWorkRepository) domain.AppRepository {
	return &AppRepositoryImpl{
		UnitOfWorkRepository: repo,
	}
}

func (a *AppRepositoryImpl) CheckAppByID(ctx context.Context, id string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM m_apps WHERE id = $1 AND deleted_at IS NULL)`
	db, err := a.GetDB()
	if err != nil {
		return false, err
	}

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Warn().Msgf(util.LogErrPrepareContext, err)
		return false, err
	}
	defer func() {
		if errClose := stmt.Close(); errClose != nil {
			log.Warn().Msgf(util.LogErrPrepareContextClose, errClose)
		}
	}()

	var exist bool
	err = stmt.QueryRowContext(ctx, id).Scan(&exist)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Warn().Msgf(util.LogErrQueryRowContextScan, err)
		}
		return false, err
	}

	if exist {
		return true, nil
	}

	return false, nil
}
