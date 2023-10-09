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

	repository2 "github.com/DueIt-Jasanya-Aturuang/doraemon/pkg/repository_old"
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

	query := regexp.QuoteMeta(`INSERT INTO m_tokens (user_id, app_id, access_token, refresh_token, remember_me) VALUES ($1, $2, $3, $4, $5)`)
	mock.ExpectBegin()
	mock.ExpectPrepare(query)
	mock.ExpectExec(query).WithArgs(
		"test", "test", "test", "test", true,
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	uowRepo := repository2.NewUnitOfWorkRepoSqlImpl(db)
	securityRepo := repository2.NewSecuritySqlRepoImpl(uowRepo)
	err = securityRepo.OpenConn(context.TODO())
	assert.NoError(t, err)

	err = securityRepo.StartTx(context.TODO(), &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	assert.NoError(t, err)

	err = securityRepo.CreateToken(context.TODO(), &model.Token{
		UserID:       "test",
		AppID:        "test",
		AcceesToken:  "test",
		RefreshToken: "test",
		RememberMe:   true,
	})
	assert.NoError(t, err)
	err = securityRepo.EndTx(err)
	assert.NoError(t, err)
	securityRepo.CloseConn()

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestSecurityGetTokenByAT(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer func() {
		err := db.Close()
		if err != nil {
			log.Err(err)
		}
	}()

	query := regexp.QuoteMeta(`SELECT id, refresh_token, app_id, remember_me FROM m_tokens WHERE access_token = $1`)
	rows := sqlmock.NewRows([]string{"id", "refresh_token", "app_id", "remember_me"})
	rows.AddRow(1, "test", "test", true)

	mock.ExpectPrepare(query)
	mock.ExpectQuery(query).WithArgs("test").WillReturnRows(rows)

	uowRepo := repository2.NewUnitOfWorkRepoSqlImpl(db)
	securityRepo := repository2.NewSecuritySqlRepoImpl(uowRepo)
	err = securityRepo.OpenConn(context.TODO())
	assert.NoError(t, err)

	token, err := securityRepo.GetTokenByAT(context.TODO(), "test")
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

	query := regexp.QuoteMeta(`UPDATE m_tokens SET refresh_token = $1, access_token = $2 WHERE id = $3`)
	mock.ExpectBegin()
	mock.ExpectPrepare(query)
	mock.ExpectExec(query).WithArgs(
		"rt", "at", 1,
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	uowRepo := repository2.NewUnitOfWorkRepoSqlImpl(db)
	securityRepo := repository2.NewSecuritySqlRepoImpl(uowRepo)
	err = securityRepo.OpenConn(context.TODO())
	assert.NoError(t, err)

	err = securityRepo.StartTx(context.TODO(), &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	assert.NoError(t, err)

	err = securityRepo.UpdateToken(context.TODO(), 1, "rt", "at")
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
		1, "user_id",
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	uowRepo := repository2.NewUnitOfWorkRepoSqlImpl(db)
	securityRepo := repository2.NewSecuritySqlRepoImpl(uowRepo)
	err = securityRepo.OpenConn(context.TODO())
	assert.NoError(t, err)

	err = securityRepo.StartTx(context.TODO(), &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	assert.NoError(t, err)

	err = securityRepo.DeleteToken(context.TODO(), 1, "user_id")
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

	uowRepo := repository2.NewUnitOfWorkRepoSqlImpl(db)
	securityRepo := repository2.NewSecuritySqlRepoImpl(uowRepo)
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
