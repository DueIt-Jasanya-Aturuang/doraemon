package _repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

type AccessRepositoryImpl struct {
	domain.UnitOfWorkRepository
}

func NewAccessRepositoryImpl(uow domain.UnitOfWorkRepository) domain.AccessRepository {
	return &AccessRepositoryImpl{
		UnitOfWorkRepository: uow,
	}
}

func (a *AccessRepositoryImpl) Create(ctx context.Context, access *domain.Access) error {
	query := `INSERT INTO m_access (role_id, user_id, app_id, access_endpoint, created_at, created_by, updated_at) 
					 VALUES ($1, $2, $3, $4, $5, $6, $7)`

	tx, err := a.GetTx()
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		log.Warn().Msgf(util.LogErrPrepareContext, err)
		return err
	}
	defer func() {
		if errClose := stmt.Close(); errClose != nil {
			log.Warn().Msgf(util.LogErrPrepareContextClose, errClose)
		}
	}()

	_, err = stmt.ExecContext(ctx,
		access.RoleId,
		access.UserId,
		access.AppId,
		access.AccessEndpoint,
		access.CreatedAt,
		access.CreatedBy,
		access.UpdatedAt,
	)

	if err != nil {
		log.Warn().Msgf(util.LogErrExecContext, err)
		return err
	}

	return nil
}

func (a *AccessRepositoryImpl) GetByUserIDAndAppID(ctx context.Context, userID string, appID string) (*domain.Access, error) {
	query := `SELECT id, app_id, user_id, role_id, access_endpoint, created_at, created_by, 
       			   updated_at, updated_by, deleted_at, deleted_by 
			FROM m_access WHERE user_id = $1 AND app_id = $2`

	conn, err := a.GetConn()
	if err != nil {
		return nil, err
	}

	stmt, err := conn.PrepareContext(ctx, query)
	if err != nil {
		log.Warn().Msgf(util.LogErrPrepareContext, err)
		return nil, err
	}
	defer func() {
		if errClose := stmt.Close(); errClose != nil {
			log.Warn().Msgf(util.LogErrPrepareContextClose, errClose)
		}
	}()

	var access domain.Access
	err = stmt.QueryRowContext(ctx, userID, appID).Scan(
		&access.ID,
		&access.AppId,
		&access.UserId,
		&access.RoleId,
		&access.AccessEndpoint,
		&access.CreatedAt,
		&access.CreatedBy,
		&access.UpdatedBy,
		&access.UpdatedBy,
		&access.DeletedAt,
		&access.DeletedBy,
	)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Warn().Msgf(util.LogErrQueryRowContextScan, err)
		}
		return nil, err
	}

	return &access, nil
}
