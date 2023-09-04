package unit

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infrastructures/repository"
)

func TestSecurityCreateToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer func() {
		err := db.Close()
		if err != nil {
			log.Err(err)
		}
	}()

	query := regexp.QuoteMeta(`INSERT INTO m_tokens (id, user_id, app_id, token) VALUES ($1, $2, $3, $4)`)
	mock.ExpectBegin()
	mock.ExpectPrepare(query)
	mock.ExpectExec(query).WithArgs(
		"test", "test", "test", "test",
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	uowRepo := repository.NewUnitOfWorkRepoSqlImpl(db)
	securityRepo := repository.NewSecuritySqlRepoImpl(uowRepo)
	err = securityRepo.OpenConn(context.TODO())
	assert.NoError(t, err)

	err = securityRepo.StartTx(context.TODO(), &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	assert.NoError(t, err)

	token, err := securityRepo.CreateToken(context.TODO(), &model.Token{
		ID:     "test",
		UserID: "test",
		AppID:  "test",
		Token:  "test",
	})
	assert.NoError(t, err)
	assert.NotNil(t, token)
	err = securityRepo.EndTx(err)
	assert.NoError(t, err)
	securityRepo.CloseConn()

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestSecurityGetTokenByIDAndUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer func() {
		err := db.Close()
		if err != nil {
			log.Err(err)
		}
	}()

	query := regexp.QuoteMeta(`SELECT id, user_id, app_id, token FROM m_tokens WHERE id = $1 AND user_id = $2`)
	rows := sqlmock.NewRows([]string{"id", "user_id", "app_id", "token"})
	rows.AddRow("test", "test", "test", "test")

	mock.ExpectPrepare(query)
	mock.ExpectQuery(query).WithArgs("test", "test").WillReturnRows(rows)

	uowRepo := repository.NewUnitOfWorkRepoSqlImpl(db)
	securityRepo := repository.NewSecuritySqlRepoImpl(uowRepo)
	err = securityRepo.OpenConn(context.TODO())
	assert.NoError(t, err)

	token, err := securityRepo.GetTokenByIDAndUserID(context.TODO(), "test", "test")
	assert.NoError(t, err)
	assert.NotNil(t, token)

	securityRepo.CloseConn()

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestSecurityUpdateToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer func() {
		err := db.Close()
		if err != nil {
			log.Err(err)
		}
	}()

	query := regexp.QuoteMeta(`UPDATE m_tokens SET token = $1, id = $2 WHERE id = $3`)
	mock.ExpectBegin()
	mock.ExpectPrepare(query)
	mock.ExpectExec(query).WithArgs(
		"test", "new_id", "old_id",
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	uowRepo := repository.NewUnitOfWorkRepoSqlImpl(db)
	securityRepo := repository.NewSecuritySqlRepoImpl(uowRepo)
	err = securityRepo.OpenConn(context.TODO())
	assert.NoError(t, err)

	err = securityRepo.StartTx(context.TODO(), &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	assert.NoError(t, err)

	err = securityRepo.UpdateToken(context.TODO(), &model.TokenUpdate{
		ID:    "new_id",
		Token: "test",
		OldID: "old_id",
	})
	assert.NoError(t, err)
	err = securityRepo.EndTx(err)
	assert.NoError(t, err)
	securityRepo.CloseConn()

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestSecurityDeleteToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer func() {
		err := db.Close()
		if err != nil {
			log.Err(err)
		}
	}()

	query := regexp.QuoteMeta(`DELETE m_tokens WHERE id = $1 AND user_id = $2`)
	mock.ExpectBegin()
	mock.ExpectPrepare(query)
	mock.ExpectExec(query).WithArgs(
		"new_id", "user_id",
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	uowRepo := repository.NewUnitOfWorkRepoSqlImpl(db)
	securityRepo := repository.NewSecuritySqlRepoImpl(uowRepo)
	err = securityRepo.OpenConn(context.TODO())
	assert.NoError(t, err)

	err = securityRepo.StartTx(context.TODO(), &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	assert.NoError(t, err)

	err = securityRepo.DeleteToken(context.TODO(), "new_id", "user_id")
	assert.NoError(t, err)
	err = securityRepo.EndTx(err)
	assert.NoError(t, err)
	securityRepo.CloseConn()

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestSecurityDeleteAllTokenByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer func() {
		err := db.Close()
		if err != nil {
			log.Err(err)
		}
	}()

	query := regexp.QuoteMeta(`DELETE m_tokens WHERE user_id = $1`)
	mock.ExpectBegin()
	mock.ExpectPrepare(query)
	mock.ExpectExec(query).WithArgs(
		"user_id",
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	uowRepo := repository.NewUnitOfWorkRepoSqlImpl(db)
	securityRepo := repository.NewSecuritySqlRepoImpl(uowRepo)
	err = securityRepo.OpenConn(context.TODO())
	assert.NoError(t, err)

	err = securityRepo.StartTx(context.TODO(), &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	assert.NoError(t, err)

	err = securityRepo.DeleteAllTokenByUserID(context.TODO(), "user_id")
	assert.NoError(t, err)
	err = securityRepo.EndTx(err)
	assert.NoError(t, err)
	securityRepo.CloseConn()

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
