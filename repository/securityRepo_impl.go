package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

type SecurityRepositoryImpl struct {
	domain.UnitOfWorkRepository
}

func NewSecurityRepositoryImpl(
	uow domain.UnitOfWorkRepository,
) domain.SecurityRepository {
	return &SecurityRepositoryImpl{
		UnitOfWorkRepository: uow,
	}
}

func (s *SecurityRepositoryImpl) Create(ctx context.Context, token *domain.Token) error {
	query := `INSERT INTO m_tokens (user_id, app_id, access_token, refresh_token, remember_me) VALUES ($1, $2, $3, $4, $5)`

	tx, err := s.GetTx()
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

	_, err = stmt.ExecContext(
		ctx,
		token.UserID,
		token.AppID,
		token.AcceesToken,
		token.RefreshToken,
		token.RememberMe,
	)
	if err != nil {
		log.Warn().Msgf(util.LogErrExecContext, err)
		return err
	}

	return nil
}

func (s *SecurityRepositoryImpl) GetByAccessToken(ctx context.Context, token string) (*domain.Token, error) {
	query := `SELECT id, user_id, refresh_token, app_id, remember_me FROM m_tokens WHERE access_token = $1`

	conn, err := s.GetConn()
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

	row := stmt.QueryRowContext(ctx, token)

	var tokenModel domain.Token
	err = row.Scan(
		&tokenModel.ID,
		&tokenModel.UserID,
		&tokenModel.RefreshToken,
		&tokenModel.AppID,
		&tokenModel.RememberMe,
	)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Warn().Msgf(util.LogErrQueryRowContextScan, err)
		}
		return nil, err
	}

	return &tokenModel, nil
}

func (s *SecurityRepositoryImpl) Update(ctx context.Context, id int, refreshToken string, accessToken string) error {
	query := `UPDATE m_tokens SET refresh_token = $1, access_token = $2 WHERE id = $3`

	tx, err := s.GetTx()
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

	_, err = stmt.ExecContext(
		ctx,
		refreshToken,
		accessToken,
		id,
	)
	if err != nil {
		log.Warn().Msgf(util.LogErrExecContext, err)
		return err
	}

	return nil
}

func (s *SecurityRepositoryImpl) Delete(ctx context.Context, id int, userID string) error {
	query := `DELETE FROM m_tokens WHERE id = $1 AND user_id = $2`

	tx, err := s.GetTx()
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

	_, err = stmt.ExecContext(
		ctx,
		id,
		userID,
	)
	if err != nil {
		log.Warn().Msgf(util.LogErrExecContext, err)
		return err
	}

	return nil
}

func (s *SecurityRepositoryImpl) DeleteAllByUserID(ctx context.Context, userID string) error {
	query := `DELETE FROM m_tokens WHERE user_id = $1`

	tx, err := s.GetTx()
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

	_, err = stmt.ExecContext(
		ctx,
		userID,
	)
	if err != nil {
		log.Warn().Msgf(util.LogErrExecContext, err)
		return err
	}

	return nil
}
