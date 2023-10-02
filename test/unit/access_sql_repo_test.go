package unit

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"

	repository2 "github.com/DueIt-Jasanya-Aturuang/doraemon/pkg/repository"
)

func TestCreateAccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer func() {
		err := db.Close()
		assert.NoError(t, err)
	}()

	query := regexp.QuoteMeta(
		`INSERT INTO m_access (role_id, user_id, app_id, access_endpoint, created_at, created_by, updated_at) 
					 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
	)

	mock.ExpectBegin()
	mock.ExpectPrepare(query)
	mock.ExpectExec(query).WithArgs(
		1, "userID", "appID", "json-marshal", 0, "userID", 0,
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	uowRepo := repository2.NewUnitOfWorkRepoSqlImpl(db)
	accessRepo := repository2.NewAccessRepoSqlImpl(uowRepo)
	err = accessRepo.OpenConn(context.TODO())
	assert.NoError(t, err)
	err = accessRepo.StartTx(context.TODO(), &sql.TxOptions{})
	assert.NoError(t, err)

	access, err := accessRepo.CreateAccess(context.TODO(), &model.Access{
		AppId:          "appID",
		UserId:         "userID",
		RoleId:         1,
		AccessEndpoint: "json-marshal",
		CreatedAt:      0,
		CreatedBy:      "userID",
		UpdatedAt:      0,
		UpdatedBy:      sql.NullString{},
		DeletedAt:      sql.NullInt64{},
		DeletedBy:      sql.NullString{},
	})
	assert.NoError(t, err)
	assert.NotNil(t, access)
	err = accessRepo.EndTx(err)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetAccessByUserIDAndAppID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer func() {
		err := db.Close()
		assert.NoError(t, err)
	}()

	query := regexp.QuoteMeta(`SELECT id, app_id, user_id, role_id, access_endpoint, created_at, created_by, 
       			   updated_at, updated_by, deleted_at, deleted_by 
			FROM m_access WHERE user_id = $1 AND app_id = $2`)

	rows := sqlmock.NewRows([]string{
		"ID", "AppID",
		"UserId", "role_id",
		"AccessEndpoint", "CreatedAt",
		"CreatedBy", "UpdatedAt",
		"UpdatedBy", "DeletedAt", "DeletedBy",
	})
	rows.AddRow(1, "appID", "userID", 1, "json-marshal", 0, "userID", 0, nil, nil, nil)

	t.Run("SUCCESS", func(t *testing.T) {
		mock.ExpectPrepare(query)
		mock.ExpectQuery(query).WithArgs("userID", "appID").WillReturnRows(rows)

		uowRepo := repository2.NewUnitOfWorkRepoSqlImpl(db)
		accessRepo := repository2.NewAccessRepoSqlImpl(uowRepo)
		err = accessRepo.OpenConn(context.TODO())
		assert.NoError(t, err)

		access, err := accessRepo.GetAccessByUserIDAndAppID(context.TODO(), "userID", "appID")
		assert.NoError(t, err)
		assert.NotNil(t, access)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("ERROR", func(t *testing.T) {
		mock.ExpectPrepare(query)
		mock.ExpectQuery(query).WithArgs("userID", "appID").WillReturnError(sql.ErrNoRows)

		uowRepo := repository2.NewUnitOfWorkRepoSqlImpl(db)
		accessRepo := repository2.NewAccessRepoSqlImpl(uowRepo)
		err = accessRepo.OpenConn(context.TODO())
		assert.NoError(t, err)

		access, err := accessRepo.GetAccessByUserIDAndAppID(context.TODO(), "userID", "appID")
		assert.Error(t, err)
		assert.Nil(t, access)
		assert.Equal(t, sql.ErrNoRows, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
