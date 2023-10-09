package user_repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

func (u *UserRepositoryImpl) CheckActivasiUser(ctx context.Context, id string) (bool, error) {
	query := "SELECT email_verified_at FROM m_users WHERE id = $1 AND deleted_at IS NULL"

	db, err := u.GetDB()
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

	return exist, nil
}

func (u *UserRepositoryImpl) Check(ctx context.Context, s string) (bool, error) {
	query := u.queryCheck(s)

	db, err := u.GetDB()
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
	err = stmt.QueryRowContext(ctx, s).Scan(&exist)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Warn().Msgf(util.LogErrQueryRowContextScan, err)
		}
		return false, err
	}

	return exist, nil
}

func (u *UserRepositoryImpl) queryCheck(s string) string {
	switch s {
	case repository.CheckUserByEmail:
		return "SELECT EXISTS(SELECT 1 FROM m_users WHERE email = $1 AND deleted_at IS NULL)"
	case repository.CheckUserByUsername:
		return "SELECT EXISTS(SELECT 1 FROM m_users WHERE username = $1 AND deleted_at IS NULL)"
	}

	return ""
}
