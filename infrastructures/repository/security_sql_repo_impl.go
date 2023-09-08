package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
	_msg "github.com/DueIt-Jasanya-Aturuang/doraemon/internal/util/msg"
)

type SecuritySqlRepoImpl struct {
	repository.UnitOfWorkSqlRepo
}

func NewSecuritySqlRepoImpl(
	uow repository.UnitOfWorkSqlRepo,
) repository.SecuritySqlRepo {
	return &SecuritySqlRepoImpl{
		UnitOfWorkSqlRepo: uow,
	}
}

func (s *SecuritySqlRepoImpl) CreateToken(ctx context.Context, token *model.Token) error {
	query := `INSERT INTO m_tokens (user_id, app_id, access_token, refresh_token, remember_me) VALUES ($1, $2, $3, $4, $5)`

	tx, err := s.GetTx()
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		log.Err(err).Msg(_msg.LogErrStartPrepareContext)
		return err
	}
	defer func() {
		if errStmt := stmt.Close(); errStmt != nil {
			log.Err(err).Msg(_msg.LogErrClosePrepareContext)
		}
	}()

	_, err = stmt.ExecContext(
		ctx,
		token.UserID,
		token.AppID,
		token.AcceesToken,
		token.RefreshToken,
		token.RememberMe,
	)
	if err != nil {
		log.Err(err).Msg(_msg.LogErrExecContext)
		return err
	}

	return nil
}

func (s *SecuritySqlRepoImpl) GetTokenByAT(
	ctx context.Context, token string,
) (*model.Token, error) {
	query := `SELECT id, refresh_token, app_id, remember_me FROM m_tokens WHERE access_token = $1`

	conn, err := s.GetConn()
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
			log.Err(err).Msg(_msg.LogErrClosePrepareContext)
		}
	}()

	row := stmt.QueryRowContext(ctx, token)

	var tokenModel model.Token
	err = row.Scan(
		&tokenModel.ID,
		&tokenModel.RefreshToken,
		&tokenModel.AppID,
		&tokenModel.RememberMe,
	)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Err(err).Msg(_msg.LogErrQueryRowContextScan)
		}
		return nil, err
	}

	return &tokenModel, nil
}

func (s *SecuritySqlRepoImpl) UpdateToken(ctx context.Context, id int, refreshToken string, accessToken string) error {
	query := `UPDATE m_tokens SET refresh_token = $1, access_token = $2 WHERE id = $3`

	tx, err := s.GetTx()
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		log.Err(err).Msg(_msg.LogErrStartPrepareContext)
		return err
	}
	defer func() {
		if errStmt := stmt.Close(); errStmt != nil {
			log.Err(err).Msg(_msg.LogErrClosePrepareContext)
		}
	}()

	_, err = stmt.ExecContext(
		ctx,
		refreshToken,
		accessToken,
		id,
	)
	if err != nil {
		log.Err(err).Msg(_msg.LogErrExecContext)
		return err
	}

	return nil
}

func (s *SecuritySqlRepoImpl) DeleteToken(ctx context.Context, id int, userID string) error {
	query := `DELETE m_tokens WHERE id = $1 AND user_id = $2`

	tx, err := s.GetTx()
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		log.Err(err).Msg(_msg.LogErrStartPrepareContext)
		return err
	}
	defer func() {
		if errStmt := stmt.Close(); errStmt != nil {
			log.Err(err).Msg(_msg.LogErrClosePrepareContext)
		}
	}()

	_, err = stmt.ExecContext(
		ctx,
		id,
		userID,
	)
	if err != nil {
		log.Err(err).Msg(_msg.LogErrExecContext)
		return err
	}

	return nil
}

func (s *SecuritySqlRepoImpl) DeleteAllTokenByUserID(ctx context.Context, userID string) error {
	query := `DELETE m_tokens WHERE user_id = $1`

	tx, err := s.GetTx()
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		log.Err(err).Msg(_msg.LogErrStartPrepareContext)
		return err
	}
	defer func() {
		if errStmt := stmt.Close(); errStmt != nil {
			log.Err(err).Msg(_msg.LogErrClosePrepareContext)
		}
	}()

	_, err = stmt.ExecContext(
		ctx,
		userID,
	)
	if err != nil {
		log.Err(err).Msg(_msg.LogErrExecContext)
		return err
	}

	return nil
}
