package _repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

type UserRepositoryImpl struct {
	domain.UnitOfWorkRepository
}

func NewUserRepositoryImpl(uow domain.UnitOfWorkRepository) domain.UserRepository {
	return &UserRepositoryImpl{
		UnitOfWorkRepository: uow,
	}
}

func (u *UserRepositoryImpl) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO m_users (id, fullname, image, username, email, password, email_verified_at, created_at, created_by, updated_at) 
					 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

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
		user.ID,
		user.FullName,
		user.Image,
		user.Username,
		user.Email,
		user.Password,
		user.EmailVerifiedAt,
		user.CreatedAt,
		user.CreatedBy,
		user.UpdatedAt,
	)
	if err != nil {
		log.Warn().Msgf(util.LogErrExecContext, err)
		return err
	}

	return nil
}

func (u *UserRepositoryImpl) UpdateActivasi(ctx context.Context, user *domain.User) error {
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

func (u *UserRepositoryImpl) UpdatePassword(ctx context.Context, user *domain.User) error {
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

func (u *UserRepositoryImpl) UpdateUsername(ctx context.Context, user *domain.User) error {
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

func (u *UserRepositoryImpl) CheckActivasiUser(ctx context.Context, id string) (bool, error) {
	query := "SELECT email_verified_at FROM m_users WHERE id = $1 AND deleted_at IS NULL"

	conn, err := u.GetConn()
	if err != nil {
		return false, err
	}

	stmt, err := conn.PrepareContext(ctx, query)
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

func (u *UserRepositoryImpl) GetByEmailOrUsername(ctx context.Context, s string) (*domain.User, error) {
	query := `SELECT id, fullname, gender, image, username, email, password, phone_number, email_verified_at, 
       				 created_at, created_by, updated_at, updated_by, deleted_at, deleted_by 
			  FROM m_users WHERE username = $1 OR email = $2 AND deleted_at IS NULL`

	conn, err := u.GetConn()
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

	var user domain.User
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

func (u *UserRepositoryImpl) Check(ctx context.Context, s string) (bool, error) {
	query := u.queryCheck(s)

	conn, err := u.GetConn()
	if err != nil {
		return false, err
	}

	stmt, err := conn.PrepareContext(ctx, query)
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

func (u *UserRepositoryImpl) Get(ctx context.Context, s string) (*domain.User, error) {
	query := u.queryGet(s)

	conn, err := u.GetConn()
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

	var user domain.User
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
	case domain.GetUserByID:
		return `SELECT id, fullname, gender, image, username, email, password, phone_number, email_verified_at, 
       				 created_at, created_by, updated_at, updated_by, deleted_at, deleted_by 
			  FROM m_users WHERE id = $1 AND deleted_at IS NULL`
	case domain.GetUserByEmail:
		return `SELECT id, fullname, gender, image, username, email, password, phone_number, email_verified_at, 
       				 created_at, created_by, updated_at, updated_by, deleted_at, deleted_by 
			  FROM m_users WHERE email = $1 AND deleted_at IS NULL`
	case domain.GetUserByUsername:
		return `SELECT id, fullname, gender, image, username, email, password, phone_number, email_verified_at, 
       				 created_at, created_by, updated_at, updated_by, deleted_at, deleted_by 
			  FROM m_users WHERE username = $1 AND deleted_at IS NULL`
	}

	return ""
}

func (u *UserRepositoryImpl) queryCheck(s string) string {
	switch s {
	case domain.CheckUserByEmail:
		return "SELECT EXISTS(SELECT 1 FROM m_users WHERE email = $1 AND deleted_at IS NULL)"
	case domain.CheckUserByUsername:
		return "SELECT EXISTS(SELECT 1 FROM m_users WHERE username = $1 AND deleted_at IS NULL)"
	}

	return ""
}
