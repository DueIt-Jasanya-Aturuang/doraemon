package _repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/util/msg"
)

type AccessRepoSqlImpl struct {
	repository.UnitOfWorkSqlRepo
}

func NewAccessRepoSqlImpl(uow repository.UnitOfWorkSqlRepo) repository.AccessSqlRepo {
	return &AccessRepoSqlImpl{
		UnitOfWorkSqlRepo: uow,
	}
}

func (a *AccessRepoSqlImpl) CreateAccess(ctx context.Context, access *model.Access) (*model.Access, error) {
	query := `INSERT INTO m_access (role_id, user_id, app_id, access_endpoint, created_at, created_by, updated_at) 
					 VALUES ($1, $2, $3, $4, $5, $6, $7)`

	tx, err := a.GetTx()
	if err != nil {
		return nil, err
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		log.Err(err).Msg(_msg.LogErrStartPrepareContext)
		return nil, err
	}
	defer func() {
		if errStmt := stmt.Close(); errStmt != nil {
			log.Err(err).Msg(_msg.LogErrClosePrepareContext)
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
		log.Err(err).Msg(_msg.LogErrExecContext)
		return nil, err
	}

	return access, nil
}

func (a *AccessRepoSqlImpl) GetAccessByUserIDAndAppID(ctx context.Context, userID string, appID string) (*model.Access, error) {
	query := `SELECT id, app_id, user_id, role_id, access_endpoint, created_at, created_by, 
       			   updated_at, updated_by, deleted_at, deleted_by 
			FROM m_access WHERE user_id = $1 AND app_id = $2`

	conn, err := a.GetConn()
	if err != nil {
		return nil, err
	}

	stmt, err := conn.PrepareContext(ctx, query)
	if err != nil {
		log.Err(err).Msg(_msg.LogErrStartPrepareContext)
		return nil, err
	}
	defer func() {
		if errStmt := stmt.Close(); errStmt != nil {
			log.Err(errStmt).Msg(_msg.LogErrClosePrepareContext)
		}
	}()

	var access model.Access
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
			log.Err(err).Msg(_msg.LogErrQueryRowContextScan)
		}
		return nil, err
	}

	return &access, nil
}
