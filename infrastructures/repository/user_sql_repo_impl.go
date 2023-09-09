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

type UserRepoSqlImpl struct {
	repository.UnitOfWorkSqlRepo
}

func NewUserRepoSqlImpl(uow repository.UnitOfWorkSqlRepo) repository.UserSqlRepo {
	return &UserRepoSqlImpl{
		UnitOfWorkSqlRepo: uow,
	}
}

func (u *UserRepoSqlImpl) CreateUser(ctx context.Context, user *model.User) error {
	query := `INSERT INTO m_users (id, fullname, image, username, email, password, email_verified_at, created_at, created_by, updated_at) 
					 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	tx, err := u.GetTx()
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
		log.Err(err).Msg(_msg.LogErrExecContext)
		return err
	}

	return nil
}

func (u *UserRepoSqlImpl) UpdateActivasiUser(ctx context.Context, user *model.User) error {
	query := `UPDATE m_users SET email_verified_at = $1, updated_at = $2, updated_by = $3 WHERE id = $4`

	tx, err := u.GetTx()
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
		user.EmailVerifiedAt,
		user.UpdatedAt,
		user.UpdatedBy,
		user.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepoSqlImpl) UpdatePasswordUser(ctx context.Context, user *model.User) error {
	query := `UPDATE m_users SET password = $1, updated_at = $2, updated_by = $3 WHERE id = $4`

	tx, err := u.GetTx()
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
		user.Password,
		user.UpdatedAt,
		user.UpdatedBy,
		user.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepoSqlImpl) CheckUserByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM m_users WHERE email = $1 AND deleted_at IS NULL)`

	conn, err := u.GetConn()
	if err != nil {
		return false, err
	}

	stmt, err := conn.PrepareContext(ctx, query)
	if err != nil {
		log.Err(err).Msg(_msg.LogErrStartPrepareContext)
		return false, err
	}
	defer func() {
		if errStmt := stmt.Close(); errStmt != nil {
			log.Err(err).Msg(_msg.LogErrClosePrepareContext)
		}
	}()

	var exists bool
	err = stmt.QueryRowContext(ctx, email).Scan(&exists)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Err(err).Msg(_msg.LogErrQueryRowContextScan)
		}
		return false, err
	}

	if exists {
		return true, nil
	}

	return false, nil
}

func (u *UserRepoSqlImpl) CheckActivasiUserByID(ctx context.Context, id string) (bool, error) {
	query := `SELECT email_verified_at FROM m_users WHERE id = $1`

	conn, err := u.GetConn()
	if err != nil {
		return false, err
	}

	stmt, err := conn.PrepareContext(ctx, query)
	if err != nil {
		log.Err(err).Msg(_msg.LogErrStartPrepareContext)
		return false, err
	}
	defer func() {
		if errStmt := stmt.Close(); errStmt != nil {
			log.Err(err).Msg(_msg.LogErrClosePrepareContext)
		}
	}()

	row := stmt.QueryRowContext(ctx, id)

	var activasi bool
	err = row.Scan(&activasi)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Err(err).Msg(_msg.LogErrQueryRowContextScan)
		}
		return false, err
	}

	return activasi, nil
}

func (u *UserRepoSqlImpl) CheckUserByUsername(ctx context.Context, username string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM m_users WHERE username = $1 AND deleted_at IS NULL)`

	conn, err := u.GetConn()
	if err != nil {
		return false, err
	}

	stmt, err := conn.PrepareContext(ctx, query)
	if err != nil {
		log.Err(err).Msg(_msg.LogErrStartPrepareContext)
		return false, err
	}
	defer func() {
		if errStmt := stmt.Close(); errStmt != nil {
			log.Err(err).Msg(_msg.LogErrClosePrepareContext)
		}
	}()

	var exists bool
	err = stmt.QueryRowContext(ctx, username).Scan(&exists)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Err(err).Msg("cannot scan query row context")
		}
		return false, err
	}

	if exists {
		return true, nil
	}

	return false, nil
}

func (u *UserRepoSqlImpl) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	query := `SELECT id, fullname, gender, image, username, email, password, phone_number, email_verified_at, 
       				 created_at, created_by, updated_at, updated_by, deleted_at, deleted_by 
			  FROM m_users WHERE id = $1`

	conn, err := u.GetConn()
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

	row := stmt.QueryRowContext(ctx, id)

	user, err := u.scanRow(row)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserRepoSqlImpl) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, fullname, gender, image, username, email, password, phone_number, email_verified_at, 
       				 created_at, created_by, updated_at, updated_by, deleted_at, deleted_by 
			  FROM m_users WHERE email = $1`

	conn, err := u.GetConn()
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

	row := stmt.QueryRowContext(ctx, email)

	user, err := u.scanRow(row)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserRepoSqlImpl) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `SELECT id, fullname, gender, image, username, email, password, phone_number, email_verified_at, 
       				 created_at, created_by, updated_at, updated_by, deleted_at, deleted_by 
			  FROM m_users WHERE username = $1`

	conn, err := u.GetConn()
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

	row := stmt.QueryRowContext(ctx, username)

	user, err := u.scanRow(row)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserRepoSqlImpl) GetUserByEmailOrUsername(ctx context.Context, emailOrUsername string) (*model.User, error) {
	query := `SELECT id, fullname, gender, image, username, email, password, phone_number, email_verified_at, 
       				 created_at, created_by, updated_at, updated_by, deleted_at, deleted_by 
			  FROM m_users WHERE username = $1 OR email = $2`

	conn, err := u.GetConn()
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

	row := stmt.QueryRowContext(ctx, emailOrUsername, emailOrUsername)

	user, err := u.scanRow(row)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserRepoSqlImpl) scanRow(row *sql.Row) (*model.User, error) {
	var user model.User

	err := row.Scan(
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
	)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Err(err).Msg(_msg.LogErrQueryRowContextScan)
		}
		return nil, err
	}

	return &user, nil
}
