package security_repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

func (s *SecurityRepositoryImpl) GetByAccessToken(ctx context.Context, token string) (*repository.Security, error) {
	query := `SELECT id, user_id, refresh_token, app_id, remember_me FROM m_tokens WHERE access_token = $1`

	db, err := s.GetDB()
	if err != nil {
		return nil, err
	}

	stmt, err := db.PrepareContext(ctx, query)
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

	var security repository.Security
	err = row.Scan(
		&security.ID,
		&security.UserID,
		&security.RefreshToken,
		&security.AppID,
		&security.RememberMe,
	)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Warn().Msgf(util.LogErrQueryRowContextScan, err)
		}
		return nil, err
	}

	return &security, nil
}
