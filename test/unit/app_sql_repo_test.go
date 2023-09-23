package unit

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	repository2 "github.com/DueIt-Jasanya-Aturuang/doraemon/pkg/_repository"
)

func TestCheckAppByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer func() {
		err := db.Close()
		assert.NoError(t, err)
	}()

	query := regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM m_apps WHERE id = $1 AND deleted_at IS NULL)`)
	rows := sqlmock.NewRows([]string{
		"exists",
	})

	t.Run("SUCCESS", func(t *testing.T) {
		rows.AddRow(false)

		mock.ExpectPrepare(query)
		mock.ExpectQuery(query).WithArgs("123").WillReturnRows(rows)

		uowRepo := repository2.NewUnitOfWorkRepoSqlImpl(db)
		appRepo := repository2.NewAppRepoSqlImpl(uowRepo)
		err = appRepo.OpenConn(context.TODO())
		assert.NoError(t, err)

		exists, err := appRepo.CheckAppByID(context.TODO(), "123")
		assert.NoError(t, err)
		assert.Equal(t, false, exists)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("SUCCESS", func(t *testing.T) {
		rows.AddRow(true)

		mock.ExpectPrepare(query)
		mock.ExpectQuery(query).WithArgs("123").WillReturnRows(rows)

		uowRepo := repository2.NewUnitOfWorkRepoSqlImpl(db)
		appRepo := repository2.NewAppRepoSqlImpl(uowRepo)
		err = appRepo.OpenConn(context.TODO())
		assert.NoError(t, err)

		exists, err := appRepo.CheckAppByID(context.TODO(), "123")
		assert.NoError(t, err)
		assert.Equal(t, true, exists)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
