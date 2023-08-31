package unit

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/infrastructures/repository"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestCheckAppByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer func() {
		err := db.Close()
		assert.NoError(t, err)
	}()

	query := regexp.QuoteMeta(`SELECT id FROM m_app WHERE id = $1`)
	rows := sqlmock.NewRows([]string{
		"exists",
	})

	t.Run("SUCCESS", func(t *testing.T) {
		rows.AddRow(false)

		mock.ExpectPrepare(query)
		mock.ExpectQuery(query).WithArgs("123").WillReturnRows(rows)

		uowRepo := repository.NewUnitOfWorkImpl(db)
		appRepo := repository.NewAppRepoImpl(uowRepo)
		err = appRepo.OpenConn(context.TODO())
		assert.NoError(t, err)

		exists, err := appRepo.CheckAppByID(context.TODO(), "123")
		assert.NoError(t, err)
		assert.Equal(t, false, exists)
	})

	t.Run("SUCCESS", func(t *testing.T) {
		rows.AddRow(true)

		mock.ExpectPrepare(query)
		mock.ExpectQuery(query).WithArgs("123").WillReturnRows(rows)

		uowRepo := repository.NewUnitOfWorkImpl(db)
		appRepo := repository.NewAppRepoImpl(uowRepo)
		err = appRepo.OpenConn(context.TODO())
		assert.NoError(t, err)

		exists, err := appRepo.CheckAppByID(context.TODO(), "123")
		assert.NoError(t, err)
		assert.Equal(t, true, exists)
	})
}
