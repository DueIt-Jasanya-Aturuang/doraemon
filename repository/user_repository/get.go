package user_repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

func (u *UserRepositoryImpl) GetByEmailOrUsername(ctx context.Context, s string) (*repository.User, error) {
	query := `SELECT id, fullname, gender, image, username, email, password, phone_number, email_verified_at, 
       				 created_at, created_by, updated_at, updated_by, deleted_at, deleted_by 
			  FROM m_users WHERE username = $1 OR email = $2 AND deleted_at IS NULL`

	db, err := u.GetDB()
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

	var user repository.User
	if err = stmt.QueryRowContext(ctx, s, s).Scan(
		&user.ID,
		&user.FullName,
		&user.Gender,
		&user.Image,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.PhoneNumber,
		&user.EmailVerifiedAt,
		&user.CreatedAt,
		&user.CreatedBy,
		&user.UpdatedAt,
		&user.UpdatedBy,
		&user.DeletedAt,
		&user.DeletedBy,
	); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Warn().Msgf(util.LogErrQueryRowContextScan, err)
		}
		return nil, err
	}

	return &user, nil
}

func (u *UserRepositoryImpl) Get(ctx context.Context, s string) (*repository.User, error) {
	query := u.queryGet(s)

	db, err := u.GetDB()
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

	var user repository.User
	if err = stmt.QueryRowContext(ctx, s).Scan(
		&user.ID,
		&user.FullName,
		&user.Gender,
		&user.Image,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.PhoneNumber,
		&user.EmailVerifiedAt,
		&user.CreatedAt,
		&user.CreatedBy,
		&user.UpdatedAt,
		&user.UpdatedBy,
		&user.DeletedAt,
		&user.DeletedBy,
	); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Warn().Msgf(util.LogErrQueryRowContextScan, err)
		}
		return nil, err
	}

	return &user, nil
}

func (u *UserRepositoryImpl) queryGet(s string) string {
	switch s {
	case repository.GetUserByID:
		return `SELECT id, fullname, gender, image, username, email, password, phone_number, email_verified_at, 
       				 created_at, created_by, updated_at, updated_by, deleted_at, deleted_by 
			  FROM m_users WHERE id = $1 AND deleted_at IS NULL`
	case repository.GetUserByEmail:
		return `SELECT id, fullname, gender, image, username, email, password, phone_number, email_verified_at, 
       				 created_at, created_by, updated_at, updated_by, deleted_at, deleted_by 
			  FROM m_users WHERE email = $1 AND deleted_at IS NULL`
	case repository.GetUserByUsername:
		return `SELECT id, fullname, gender, image, username, email, password, phone_number, email_verified_at, 
       				 created_at, created_by, updated_at, updated_by, deleted_at, deleted_by 
			  FROM m_users WHERE username = $1 AND deleted_at IS NULL`
	}

	return ""
}
