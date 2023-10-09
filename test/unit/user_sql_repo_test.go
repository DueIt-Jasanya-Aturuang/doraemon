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

func TestCreateUserRepo(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer func() {
		err := db.Close()
		if err != nil {
			log.Err(err)
		}
	}()

	query := regexp.QuoteMeta(`INSERT INTO m_users (id, fullname, image, username, email, password, email_verified_at, created_at, created_by, updated_at) 
					 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`)

	mock.ExpectBegin()
	mock.ExpectPrepare(query)
	mock.ExpectExec(query).WithArgs(
		"test", "test", "default-male.png", "rmaa", "ibanrama29@gmail.com", "123", false, 0, "test", 0,
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	uowRepo := repository2.NewUnitOfWorkRepoSqlImpl(db)
	userRepo := repository2.NewUserRepoSqlImpl(uowRepo)
	err = userRepo.OpenConn(context.TODO())
	assert.NoError(t, err)
	err = userRepo.StartTx(context.TODO(), &sql.TxOptions{})
	assert.NoError(t, err)

	err = userRepo.CreateUser(context.TODO(), &model.User{
		ID:              "test",
		FullName:        "test",
		Gender:          "undefined",
		Image:           "default-male.png",
		Username:        "rmaa",
		Email:           "ibanrama29@gmail.com",
		Password:        "123",
		PhoneNumber:     sql.NullString{},
		EmailVerifiedAt: false,
		CreatedAt:       0,
		CreatedBy:       "test",
		UpdatedAt:       0,
		UpdatedBy:       sql.NullString{},
		DeletedAt:       sql.NullInt64{},
		DeletedBy:       sql.NullString{},
	})
	assert.NoError(t, err)

	err = userRepo.EndTx(err)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdateActivatiUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer func() {
		err := db.Close()
		if err != nil {
			log.Err(err)
		}
	}()

	query := regexp.QuoteMeta(`UPDATE m_users SET email_verified_at = $1, updated_at = $2, updated_by = $3 WHERE id = $4`)

	mock.ExpectBegin()
	mock.ExpectPrepare(query)
	mock.ExpectExec(query).WithArgs(
		true, 0, "test", "test",
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	uowRepo := repository2.NewUnitOfWorkRepoSqlImpl(db)
	userRepo := repository2.NewUserRepoSqlImpl(uowRepo)
	err = userRepo.OpenConn(context.TODO())
	assert.NoError(t, err)
	err = userRepo.StartTx(context.TODO(), &sql.TxOptions{})
	assert.NoError(t, err)

	err = userRepo.UpdateActivasiUser(context.TODO(), &model.User{
		ID:              "test",
		EmailVerifiedAt: true,
		UpdatedAt:       0,
		UpdatedBy:       sql.NullString{String: "test", Valid: true},
	})
	assert.NoError(t, err)
	err = userRepo.EndTx(err)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdatePasswordUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer func() {
		err := db.Close()
		if err != nil {
			log.Err(err)
		}
	}()

	query := regexp.QuoteMeta(`UPDATE m_users SET password = $1, updated_at = $2, updated_by = $3 WHERE id = $4`)

	mock.ExpectBegin()
	mock.ExpectPrepare(query)
	mock.ExpectExec(query).WithArgs(
		"true", 0, "test", "test",
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	uowRepo := repository2.NewUnitOfWorkRepoSqlImpl(db)
	userRepo := repository2.NewUserRepoSqlImpl(uowRepo)
	err = userRepo.OpenConn(context.TODO())
	assert.NoError(t, err)
	err = userRepo.StartTx(context.TODO(), &sql.TxOptions{})
	assert.NoError(t, err)

	err = userRepo.UpdatePasswordUser(context.TODO(), &model.User{
		ID:        "test",
		Password:  "true",
		UpdatedAt: 0,
		UpdatedBy: sql.NullString{String: "test", Valid: true},
	})
	assert.NoError(t, err)
	err = userRepo.EndTx(err)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCheckEmailUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer func() {
		err := db.Close()
		assert.NoError(t, err)
	}()

	query := regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM m_users WHERE email = $1 AND deleted_at IS NULL)`)
	rows := sqlmock.NewRows([]string{
		"exists",
	})

	t.Run("SUCCESS", func(t *testing.T) {
		rows.AddRow(false)

		mock.ExpectPrepare(query)
		mock.ExpectQuery(query).WithArgs("ibanrama29@gmail.com").WillReturnRows(rows)

		uowRepo := repository2.NewUnitOfWorkRepoSqlImpl(db)
		userRepo := repository2.NewUserRepoSqlImpl(uowRepo)
		err = userRepo.OpenConn(context.TODO())
		assert.NoError(t, err)

		exists, err := userRepo.CheckUserByEmail(context.TODO(), "ibanrama29@gmail.com")
		assert.NoError(t, err)
		assert.Equal(t, false, exists)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("ERROR", func(t *testing.T) {
		rows.AddRow(true)

		mock.ExpectPrepare(query)
		mock.ExpectQuery(query).WithArgs("ibanrama29@gmail.com").WillReturnRows(rows)

		uowRepo := repository2.NewUnitOfWorkRepoSqlImpl(db)
		userRepo := repository2.NewUserRepoSqlImpl(uowRepo)
		err = userRepo.OpenConn(context.TODO())
		assert.NoError(t, err)

		exists, err := userRepo.CheckUserByEmail(context.TODO(), "ibanrama29@gmail.com")
		assert.NoError(t, err)
		assert.Equal(t, true, exists)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestCheckUsernameUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer func() {
		err := db.Close()
		assert.NoError(t, err)
	}()

	query := regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM m_users WHERE username = $1 AND deleted_at IS NULL)`)
	rows := sqlmock.NewRows([]string{
		"exists",
	})

	t.Run("SUCCESS", func(t *testing.T) {
		rows.AddRow(false)

		mock.ExpectPrepare(query)
		mock.ExpectQuery(query).WithArgs("rmaa").WillReturnRows(rows)

		uowRepo := repository2.NewUnitOfWorkRepoSqlImpl(db)
		userRepo := repository2.NewUserRepoSqlImpl(uowRepo)
		err = userRepo.OpenConn(context.TODO())
		assert.NoError(t, err)

		exists, err := userRepo.CheckUserByUsername(context.TODO(), "rmaa")
		assert.NoError(t, err)
		assert.Equal(t, false, exists)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("ERROR", func(t *testing.T) {
		rows.AddRow(true)

		mock.ExpectPrepare(query)
		mock.ExpectQuery(query).WithArgs("rmaa").WillReturnRows(rows)

		uowRepo := repository2.NewUnitOfWorkRepoSqlImpl(db)
		userRepo := repository2.NewUserRepoSqlImpl(uowRepo)
		err = userRepo.OpenConn(context.TODO())
		assert.NoError(t, err)

		exists, err := userRepo.CheckUserByUsername(context.TODO(), "rmaa")
		assert.NoError(t, err)
		assert.Equal(t, true, exists)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestGetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer func() {
		err := db.Close()
		assert.NoError(t, err)
	}()

	query := regexp.QuoteMeta(`SELECT id, fullname, gender, image, username, email, password, phone_number, email_verified_at, 
       				 created_at, created_by, updated_at, updated_by, deleted_at, deleted_by 
			  FROM m_users WHERE id = $1`)

	rows := sqlmock.NewRows([]string{
		"id", "fullname", "gender", "image", "username", "email", "password", "phone_number",
		"email_verified_at", "created_at", "created_by", "updated_at",
		"updated_by", "deleted_at", "deleted_by",
	})

	t.Run("SUCCESS", func(t *testing.T) {
		rows.AddRow("userID", "syaiban ahmad ramadhan", "undefined", "default-male.png", "rama",
			"ibanrama29@gmail.com", "123456", nil, false, 0, "userID", 0, nil, nil, nil,
		)

		mock.ExpectPrepare(query)
		mock.ExpectQuery(query).WithArgs("userID").WillReturnRows(rows)

		uowRepo := repository2.NewUnitOfWorkRepoSqlImpl(db)
		userRepo := repository2.NewUserRepoSqlImpl(uowRepo)
		err = userRepo.OpenConn(context.TODO())
		assert.NoError(t, err)

		user, err := userRepo.GetUserByID(context.TODO(), "userID")
		assert.NoError(t, err)
		assert.NotNil(t, user)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("ERROR", func(t *testing.T) {
		mock.ExpectPrepare(query)
		mock.ExpectQuery(query).WithArgs("rmaa").WillReturnError(sql.ErrNoRows)

		uowRepo := repository2.NewUnitOfWorkRepoSqlImpl(db)
		userRepo := repository2.NewUserRepoSqlImpl(uowRepo)
		err = userRepo.OpenConn(context.TODO())
		assert.NoError(t, err)

		user, err := userRepo.GetUserByID(context.TODO(), "rmaa")
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.Nil(t, user)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestGetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer func() {
		err := db.Close()
		assert.NoError(t, err)
	}()

	query := regexp.QuoteMeta(`SELECT id, fullname, gender, image, username, email, password, phone_number, email_verified_at, 
       				 created_at, created_by, updated_at, updated_by, deleted_at, deleted_by 
			  FROM m_users WHERE email = $1`)

	rows := sqlmock.NewRows([]string{
		"id", "fullname", "gender", "image", "username", "email", "password", "phone_number",
		"email_verified_at", "created_at", "created_by", "updated_at",
		"updated_by", "deleted_at", "deleted_by",
	})

	t.Run("SUCCESS", func(t *testing.T) {
		rows.AddRow("userID", "syaiban ahmad ramadhan", "undefined", "default-male.png", "rama",
			"ibanrama29@gmail.com", "123456", nil, false, 0, "userID", 0, nil, nil, nil,
		)

		mock.ExpectPrepare(query)
		mock.ExpectQuery(query).WithArgs("ibanrama95@gmail.com").WillReturnRows(rows)

		uowRepo := repository2.NewUnitOfWorkRepoSqlImpl(db)
		userRepo := repository2.NewUserRepoSqlImpl(uowRepo)
		err = userRepo.OpenConn(context.TODO())
		assert.NoError(t, err)

		user, err := userRepo.GetUserByEmail(context.TODO(), "ibanrama95@gmail.com")
		assert.NoError(t, err)
		assert.NotNil(t, user)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("ERROR", func(t *testing.T) {
		mock.ExpectPrepare(query)
		mock.ExpectQuery(query).WithArgs("nil").WillReturnError(sql.ErrNoRows)

		uowRepo := repository2.NewUnitOfWorkRepoSqlImpl(db)
		userRepo := repository2.NewUserRepoSqlImpl(uowRepo)
		err = userRepo.OpenConn(context.TODO())
		assert.NoError(t, err)

		user, err := userRepo.GetUserByEmail(context.TODO(), "nil")
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.Nil(t, user)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestGetUserByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer func() {
		err := db.Close()
		assert.NoError(t, err)
	}()

	query := regexp.QuoteMeta(`SELECT id, fullname, gender, image, username, email, password, phone_number, email_verified_at, 
       				 created_at, created_by, updated_at, updated_by, deleted_at, deleted_by 
			  FROM m_users WHERE username = $1`)

	rows := sqlmock.NewRows([]string{
		"id", "fullname", "gender", "image", "username", "email", "password", "phone_number",
		"email_verified_at", "created_at", "created_by", "updated_at",
		"updated_by", "deleted_at", "deleted_by",
	})

	t.Run("SUCCESS", func(t *testing.T) {
		rows.AddRow("userID", "syaiban ahmad ramadhan", "undefined", "default-male.png", "rama",
			"ibanrama29@gmail.com", "123456", nil, false, 0, "userID", 0, nil, nil, nil,
		)

		mock.ExpectPrepare(query)
		mock.ExpectQuery(query).WithArgs("rama").WillReturnRows(rows)

		uowRepo := repository2.NewUnitOfWorkRepoSqlImpl(db)
		userRepo := repository2.NewUserRepoSqlImpl(uowRepo)
		err = userRepo.OpenConn(context.TODO())
		assert.NoError(t, err)

		user, err := userRepo.GetUserByUsername(context.TODO(), "rama")
		assert.NoError(t, err)
		assert.NotNil(t, user)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("ERROR", func(t *testing.T) {
		mock.ExpectPrepare(query)
		mock.ExpectQuery(query).WithArgs("nil").WillReturnError(sql.ErrNoRows)

		uowRepo := repository2.NewUnitOfWorkRepoSqlImpl(db)
		userRepo := repository2.NewUserRepoSqlImpl(uowRepo)
		err = userRepo.OpenConn(context.TODO())
		assert.NoError(t, err)

		user, err := userRepo.GetUserByUsername(context.TODO(), "nil")
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.Nil(t, user)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestGetUserByEmailOrUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer func() {
		err := db.Close()
		assert.NoError(t, err)
	}()

	query := regexp.QuoteMeta(`SELECT id, fullname, gender, image, username, email, password, phone_number, email_verified_at, 
       				 created_at, created_by, updated_at, updated_by, deleted_at, deleted_by 
			  FROM m_users WHERE username = $1 OR email = $2`)

	rows := sqlmock.NewRows([]string{
		"id", "fullname", "gender", "image", "username", "email", "password", "phone_number",
		"email_verified_at", "created_at", "created_by", "updated_at",
		"updated_by", "deleted_at", "deleted_by",
	})

	t.Run("SUCCESS", func(t *testing.T) {
		rows.AddRow("userID", "syaiban ahmad ramadhan", "undefined", "default-male.png", "rama",
			"ibanrama29@gmail.com", "123456", nil, false, 0, "userID", 0, nil, nil, nil,
		)

		mock.ExpectPrepare(query)
		mock.ExpectQuery(query).WithArgs("rama", "rama").WillReturnRows(rows)

		uowRepo := repository2.NewUnitOfWorkRepoSqlImpl(db)
		userRepo := repository2.NewUserRepoSqlImpl(uowRepo)
		err = userRepo.OpenConn(context.TODO())
		assert.NoError(t, err)

		user, err := userRepo.GetUserByEmailOrUsername(context.TODO(), "rama")
		assert.NoError(t, err)
		assert.NotNil(t, user)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("ERROR", func(t *testing.T) {
		mock.ExpectPrepare(query)
		mock.ExpectQuery(query).WithArgs("nil", "nil").WillReturnError(sql.ErrNoRows)

		uowRepo := repository2.NewUnitOfWorkRepoSqlImpl(db)
		userRepo := repository2.NewUserRepoSqlImpl(uowRepo)
		err = userRepo.OpenConn(context.TODO())
		assert.NoError(t, err)

		user, err := userRepo.GetUserByEmailOrUsername(context.TODO(), "nil")
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
		assert.Nil(t, user)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestCheckActivasiUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer func() {
		err := db.Close()
		assert.NoError(t, err)
	}()

	query := regexp.QuoteMeta(`SELECT email_verified_at FROM m_users WHERE id = $1`)

	rows := sqlmock.NewRows([]string{
		"email_verified_at",
	})

	t.Run("SUCCESS_true", func(t *testing.T) {
		rows.AddRow(true)

		mock.ExpectPrepare(query)
		mock.ExpectQuery(query).WithArgs("rama").WillReturnRows(rows)

		uowRepo := repository2.NewUnitOfWorkRepoSqlImpl(db)
		userRepo := repository2.NewUserRepoSqlImpl(uowRepo)
		err = userRepo.OpenConn(context.TODO())
		assert.NoError(t, err)

		activasi, err := userRepo.CheckActivasiUserByID(context.TODO(), "rama")
		assert.NoError(t, err)
		assert.Equal(t, true, activasi)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("SUCCESS_false", func(t *testing.T) {
		rows.AddRow(false)

		mock.ExpectPrepare(query)
		mock.ExpectQuery(query).WithArgs("rama").WillReturnRows(rows)

		uowRepo := repository2.NewUnitOfWorkRepoSqlImpl(db)
		userRepo := repository2.NewUserRepoSqlImpl(uowRepo)
		err = userRepo.OpenConn(context.TODO())
		assert.NoError(t, err)

		activasi, err := userRepo.CheckActivasiUserByID(context.TODO(), "rama")
		assert.NoError(t, err)
		assert.Equal(t, false, activasi)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
