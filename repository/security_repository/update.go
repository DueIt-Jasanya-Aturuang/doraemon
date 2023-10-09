package security_repository

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

func (s *SecurityRepositoryImpl) Update(ctx context.Context, param repository.UpdateSecurityParams) error {
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
		param.RefreshToken,
		param.AccessToken,
		param.ID,
	)
	if err != nil {
		log.Warn().Msgf(util.LogErrExecContext, err)
		return err
	}

	return nil
}
