package repository

import (
	"context"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
	"github.com/rs/zerolog/log"
)

type UserRepoImpl struct {
	repository.UnitOfWorkRepo
}

func NewUserRepoImpl(uow repository.UnitOfWorkRepo) repository.UserRepo {
	return &UserRepoImpl{
		UnitOfWorkRepo: uow,
	}
}

func (u UserRepoImpl) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	query := `INSERT INTO m_users (id, fullname, image, username, email, password, email_verified_at, created_at, created_by, updated_at) 
					 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	tx, err := u.GetTx()
	if err != nil {
		return nil, err
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		log.Err(err).Msg("failed to start prepared context")
		return nil, err
	}
	defer func() {
		if errStmt := stmt.Close(); errStmt != nil {
			log.Err(err).Msg("failed to close prepared context")
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
		log.Err(err).Msg("failed to query row context prepared statement")
		return nil, err
	}

	user.Gender = "undefined"
	return user, nil
}

func (u UserRepoImpl) CheckEmailUser(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM m_users WHERE email = $1 AND deleted_at IS NULL)`

	conn, err := u.GetConn()
	if err != nil {
		return false, err
	}

	stmt, err := conn.PrepareContext(ctx, query)
	if err != nil {
		log.Err(err).Msg("failed to start prepared context")
		return false, err
	}
	defer func() {
		if errStmt := stmt.Close(); errStmt != nil {
			log.Err(err).Msg("failed to close prepared context")
		}
	}()

	var exists bool
	err = stmt.QueryRowContext(ctx, email).Scan(&exists)
	if err != nil {
		log.Err(err).Msg("failed to query row context prepared statement")
		return false, err
	}

	if exists {
		return true, nil
	}

	return false, nil
}

func (u UserRepoImpl) CheckUsernameUser(ctx context.Context, username string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM m_users WHERE username = $1 AND deleted_at IS NULL)`

	conn, err := u.GetConn()
	if err != nil {
		return false, err
	}

	stmt, err := conn.PrepareContext(ctx, query)
	if err != nil {
		log.Err(err).Msg("failed to start prepared context")
		return false, err
	}

	var exists bool
	err = stmt.QueryRowContext(ctx, username).Scan(&exists)
	if err != nil {
		log.Err(err).Msg("failed to query row context prepared statement")
		return false, err
	}

	if exists {
		return true, nil
	}

	return false, nil
}
