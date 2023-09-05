package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
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
	query := `INSERT INTO m_tokens (id, user_id, app_id, token, remember_me) VALUES ($1, $2, $3, $4, $5)`

	tx, err := s.GetTx()
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		log.Err(err).Msg("failed to start prepared context")
		return err
	}
	defer func() {
		if errStmt := stmt.Close(); errStmt != nil {
			log.Err(err).Msg("failed to close prepared context")
		}
	}()

	_, err = stmt.ExecContext(
		ctx,
		token.ID,
		token.UserID,
		token.AppID,
		token.Token,
		token.RememberMe,
	)
	if err != nil {
		log.Err(err).Msg("failed to query row context prepared statement")
		return err
	}

	return nil
}

func (s *SecuritySqlRepoImpl) GetTokenByIDAndUserID(
	ctx context.Context, tokenID string, userID string,
) (*model.Token, error) {
	query := `SELECT id, user_id, app_id, token FROM m_tokens WHERE id = $1 AND user_id = $2`

	conn, err := s.GetConn()
	if err != nil {
		return nil, err
	}

	stmt, err := conn.PrepareContext(ctx, query)
	if err != nil {
		log.Err(err).Msg("failed to start prepared context")
		return nil, err
	}
	defer func() {
		if errStmt := stmt.Close(); errStmt != nil {
			log.Err(err).Msg("failed to close prepared context")
		}
	}()

	row := stmt.QueryRowContext(ctx, tokenID, userID)

	var token model.Token
	err = row.Scan(
		&token.ID,
		&token.UserID,
		&token.AppID,
		&token.Token,
	)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Err(err).Msg("cannot scan query row context")
		}
		return nil, err
	}

	return &token, nil
}

func (s *SecuritySqlRepoImpl) UpdateToken(ctx context.Context, token *model.TokenUpdate) error {
	query := `UPDATE m_tokens SET token = $1, id = $2 WHERE id = $3`

	tx, err := s.GetTx()
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		log.Err(err).Msg("failed to start prepared context")
		return err
	}
	defer func() {
		if errStmt := stmt.Close(); errStmt != nil {
			log.Err(err).Msg("failed to close prepared context")
		}
	}()

	_, err = stmt.ExecContext(
		ctx,
		token.Token,
		token.ID,
		token.OldID,
	)
	if err != nil {
		log.Err(err).Msg("failed to query row context prepared statement")
		return err
	}

	return nil
}

func (s *SecuritySqlRepoImpl) DeleteToken(ctx context.Context, tokenID string, userID string) error {
	query := `DELETE m_tokens WHERE id = $1 AND user_id = $2`

	tx, err := s.GetTx()
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		log.Err(err).Msg("failed to start prepared context")
		return err
	}
	defer func() {
		if errStmt := stmt.Close(); errStmt != nil {
			log.Err(err).Msg("failed to close prepared context")
		}
	}()

	_, err = stmt.ExecContext(
		ctx,
		tokenID,
		userID,
	)
	if err != nil {
		log.Err(err).Msg("failed to query row context prepared statement")
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
		log.Err(err).Msg("failed to start prepared context")
		return err
	}
	defer func() {
		if errStmt := stmt.Close(); errStmt != nil {
			log.Err(err).Msg("failed to close prepared context")
		}
	}()

	_, err = stmt.ExecContext(
		ctx,
		userID,
	)
	if err != nil {
		log.Err(err).Msg("failed to query row context prepared statement")
		return err
	}

	return nil
}
