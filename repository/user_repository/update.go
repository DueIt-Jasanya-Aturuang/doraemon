package user_repository

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

func (u *UserRepositoryImpl) UpdateActivasi(ctx context.Context, user *repository.User) error {
	query := `UPDATE m_users SET email_verified_at = $1, updated_at = $2, updated_by = $3 WHERE id = $4`

	tx, err := u.GetTx()
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
		user.EmailVerifiedAt,
		user.UpdatedAt,
		user.UpdatedBy,
		user.ID,
	)
	if err != nil {
		log.Warn().Msgf(util.LogErrExecContext, err)
		return err
	}

	return nil
}

func (u *UserRepositoryImpl) UpdatePassword(ctx context.Context, user *repository.User) error {
	query := `UPDATE m_users SET password = $1, updated_at = $2, updated_by = $3 WHERE id = $4`

	tx, err := u.GetTx()
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
		user.Password,
		user.UpdatedAt,
		user.UpdatedBy,
		user.ID,
	)
	if err != nil {
		log.Warn().Msgf(util.LogErrExecContext, err)
		return err
	}

	return nil
}

func (u *UserRepositoryImpl) UpdateUsername(ctx context.Context, user *repository.User) error {
	query := `UPDATE m_users SET username = $1, updated_at = $2, updated_by = $3 WHERE id = $4`

	tx, err := u.GetTx()
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
		user.Username,
		user.UpdatedAt,
		user.UpdatedBy,
		user.ID,
	)
	if err != nil {
		log.Warn().Msgf(util.LogErrExecContext, err)
		return err
	}

	return nil
}
