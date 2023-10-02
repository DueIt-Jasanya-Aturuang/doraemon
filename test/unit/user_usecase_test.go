package unit

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/mocks"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase/converter"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase/helper"
)

func TestUserUsecaseResetPassword(t *testing.T) {
	userRepo := &mocks.FakeUserSqlRepo{}
	redisClient, _ := redismock.NewClientMock()
	redisImpl := &infra.RedisImpl{Client: redisClient}

	userUsecase := usecase.NewUserUsecaseImpl(userRepo, redisImpl)

	req := &dto.ResetPasswordReq{
		OldPassword: "old12345",
		Password:    "new12345",
		RePassword:  "new12345",
		UserID:      "userID_1",
	}

	user := &model.User{
		ID:              "userID_1",
		FullName:        "ibanrama",
		Gender:          "undefined",
		Image:           ".png",
		Username:        "rama",
		Email:           "ibanrama29@gmail.com",
		Password:        "$2a$14$bUlfMzLFLKrILX3uNBkEieHMt/0boFfBvzpyDN4ai.UytAaTwmP1u",
		PhoneNumber:     sql.NullString{},
		EmailVerifiedAt: true,
		CreatedAt:       0,
		CreatedBy:       "userID_1",
		UpdatedAt:       0,
		UpdatedBy:       sql.NullString{},
		DeletedAt:       sql.NullInt64{},
		DeletedBy:       sql.NullString{},
	}

	t.Run("SUCCESS", func(t *testing.T) {
		_ = userRepo.OpenConn(context.TODO())
		defer userRepo.CloseConn()

		_, _ = userRepo.GetUserByID(context.TODO(), req.UserID)
		userRepo.GetUserByIDReturns(user, nil)

		_ = userRepo.StartTx(context.TODO(), &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
			ReadOnly:  false,
		})
		defer func() {
			_ = userRepo.EndTx(nil)
		}()
		userRepo.EndTxReturns(nil)

		userConv := converter.ChangePasswordReqToModel("passwordHash", user.ID)
		_ = userRepo.UpdatePasswordUser(context.TODO(), userConv)
		userRepo.UpdatePasswordUserReturns(nil)

		err := userUsecase.ChangePassword(context.TODO(), req)
		assert.NoError(t, err)
	})

	t.Run("ERROR_invalid-password", func(t *testing.T) {
		req.OldPassword = "invalid"
		_ = userRepo.OpenConn(context.TODO())
		defer userRepo.CloseConn()

		_, _ = userRepo.GetUserByID(context.TODO(), req.UserID)
		userRepo.GetUserByIDReturns(user, nil)

		_ = userRepo.StartTx(context.TODO(), &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
			ReadOnly:  false,
		})
		defer func() {
			_ = userRepo.EndTx(errors.New("err"))
		}()
		userRepo.EndTxReturns(nil)

		userConv := converter.ChangePasswordReqToModel("passwordHash", user.ID)
		_ = userRepo.UpdatePasswordUser(context.TODO(), userConv)
		userRepo.UpdatePasswordUserReturns(nil)

		err := userUsecase.ChangePassword(context.TODO(), req)

		assert.Error(t, err)
		var errHTTP *model.ErrResponseHTTP
		assert.Equal(t, true, errors.As(err, &errHTTP))
		assert.Equal(t, 400, errHTTP.Code)
	})

	t.Run("ERROR_invalid-user", func(t *testing.T) {
		req.OldPassword = "invalid"
		_ = userRepo.OpenConn(context.TODO())
		defer userRepo.CloseConn()

		_, _ = userRepo.GetUserByID(context.TODO(), "nil")
		userRepo.GetUserByIDReturns(nil, sql.ErrNoRows)

		err := userUsecase.ChangePassword(context.TODO(), req)
		assert.Error(t, err)
		var errHTTP *model.ErrResponseHTTP
		assert.Equal(t, true, errors.As(err, &errHTTP))
		assert.Equal(t, 401, errHTTP.Code)
	})
}

func TestUserUsecaseForgottenPassword(t *testing.T) {
	infra.EnvInit()
	userRepo := &mocks.FakeUserSqlRepo{}
	redisClient, mock := redismock.NewClientMock()
	redisImpl := &infra.RedisImpl{Client: redisClient}

	userUsecase := usecase.NewUserUsecaseImpl(userRepo, redisImpl)

	req := &dto.ForgottenPasswordReq{
		Email: "ibanrama29@gmail.com",
	}

	user := &model.User{
		ID:              "userID_1",
		FullName:        "ibanrama",
		Gender:          "undefined",
		Image:           ".png",
		Username:        "rama",
		Email:           "ibanrama29@gmail.com",
		Password:        "$2a$14$bUlfMzLFLKrILX3uNBkEieHMt/0boFfBvzpyDN4ai.UytAaTwmP1u",
		PhoneNumber:     sql.NullString{},
		EmailVerifiedAt: true,
		CreatedAt:       0,
		CreatedBy:       "userID_1",
		UpdatedAt:       0,
		UpdatedBy:       sql.NullString{},
		DeletedAt:       sql.NullInt64{},
		DeletedBy:       sql.NullString{},
	}

	t.Run("SUCCESS", func(t *testing.T) {
		_ = userRepo.OpenConn(context.TODO())
		defer userRepo.CloseConn()

		_, _ = userRepo.GetUserByEmail(context.TODO(), req.Email)
		userRepo.GetUserByEmailReturns(user, nil)

		mock.ExpectSet("forgot-password-link:"+req.Email, user.ID, 5*time.Minute).SetVal(user.ID)

		link, err := userUsecase.ForgottenPassword(context.TODO(), req)
		assert.NotEqual(t, " ", link)
		assert.NoError(t, err)
		t.Log(link)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("ERROR_invalid-email", func(t *testing.T) {
		_ = userRepo.OpenConn(context.TODO())
		defer userRepo.CloseConn()

		_, _ = userRepo.GetUserByEmail(context.TODO(), "nil")
		userRepo.GetUserByEmailReturns(nil, sql.ErrNoRows)

		link, err := userUsecase.ForgottenPassword(context.TODO(), req)
		t.Log(link)
		assert.NotEqual(t, " ", link)
		assert.Error(t, err)
		var errHTTP *model.ErrResponseHTTP
		assert.Equal(t, true, errors.As(err, &errHTTP))
		assert.Equal(t, 404, errHTTP.Code)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestUserUsecaseResetForgottenPassword(t *testing.T) {
	infra.EnvInit()
	userRepo := &mocks.FakeUserSqlRepo{}
	redisClient, mock := redismock.NewClientMock()
	redisImpl := &infra.RedisImpl{Client: redisClient}

	userUsecase := usecase.NewUserUsecaseImpl(userRepo, redisImpl)

	var jwtModel *model.Jwt
	jwtModel = jwtModel.ForgotPasswordTokenDefault("userID_1")
	forgotPasswordTokenRes, err := helper.GenerateJwtHS256(jwtModel)
	assert.NoError(t, err)

	req := &dto.ResetForgottenPasswordReq{
		Email:      "ibanrama29@gmail.com",
		Token:      forgotPasswordTokenRes,
		Password:   "new12345",
		RePassword: "new12345",
	}

	user := &model.User{
		ID:              "userID_1",
		FullName:        "ibanrama",
		Gender:          "undefined",
		Image:           ".png",
		Username:        "rama",
		Email:           "ibanrama29@gmail.com",
		Password:        "$2a$14$bUlfMzLFLKrILX3uNBkEieHMt/0boFfBvzpyDN4ai.UytAaTwmP1u",
		PhoneNumber:     sql.NullString{},
		EmailVerifiedAt: true,
		CreatedAt:       0,
		CreatedBy:       "userID_1",
		UpdatedAt:       0,
		UpdatedBy:       sql.NullString{},
		DeletedAt:       sql.NullInt64{},
		DeletedBy:       sql.NullString{},
	}

	t.Run("SUCCESS", func(t *testing.T) {
		_ = userRepo.OpenConn(context.TODO())
		defer userRepo.CloseConn()

		_ = userRepo.StartTx(context.TODO(), &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
			ReadOnly:  false,
		})
		userRepo.StartTxReturns(nil)
		defer func() {
			_ = userRepo.EndTx(nil)
			userRepo.EndTxReturns(nil)
		}()

		mock.ExpectGet("forgot-password-link:" + req.Email).SetVal(user.ID)

		userConv := converter.ChangePasswordReqToModel("passwordHash", "userID_1")

		_ = userRepo.UpdatePasswordUser(context.TODO(), userConv)
		userRepo.UpdatePasswordUserReturns(nil)

		err := userUsecase.ResetForgottenPassword(context.TODO(), req)
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("ERROR_invalid-token", func(t *testing.T) {
		req.Token = "invalid token"
		_ = userRepo.OpenConn(context.TODO())
		defer userRepo.CloseConn()

		_ = userRepo.StartTx(context.TODO(), &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
			ReadOnly:  false,
		})
		userRepo.StartTxReturns(nil)
		defer func() {
			_ = userRepo.EndTx(nil)
			userRepo.EndTxReturns(nil)
		}()

		userConv := converter.ChangePasswordReqToModel("passwordHash", "userID_1")

		_ = userRepo.UpdatePasswordUser(context.TODO(), userConv)
		userRepo.UpdatePasswordUserReturns(nil)

		err := userUsecase.ResetForgottenPassword(context.TODO(), req)
		assert.Error(t, err)
		var errHTTP *model.ErrResponseHTTP
		assert.Equal(t, true, errors.As(err, &errHTTP))
		assert.Equal(t, 401, errHTTP.Code)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
		mock.ClearExpect()
	})

	t.Run("ERROR_null-in-redis", func(t *testing.T) {
		req.Token = forgotPasswordTokenRes
		mock.ClearExpect()
		_ = userRepo.OpenConn(context.TODO())
		defer userRepo.CloseConn()

		_ = userRepo.StartTx(context.TODO(), &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
			ReadOnly:  false,
		})
		userRepo.StartTxReturns(nil)
		defer func() {
			_ = userRepo.EndTx(nil)
			userRepo.EndTxReturns(nil)
		}()

		mock.ExpectGet("forgot-password-link:" + req.Email).SetVal("")

		err = userUsecase.ResetForgottenPassword(context.TODO(), req)
		assert.Error(t, err)
		var errHTTP *model.ErrResponseHTTP
		assert.Equal(t, true, errors.As(err, &errHTTP))
		assert.Equal(t, 401, errHTTP.Code)
		errMock := mock.ExpectationsWereMet()
		assert.NoError(t, errMock)
	})
}

func TestUserUsecaseActivasiAccount(t *testing.T) {
	infra.EnvInit()
	userRepo := &mocks.FakeUserSqlRepo{}
	redisClient, mock := redismock.NewClientMock()
	redisImpl := &infra.RedisImpl{Client: redisClient}

	userUsecase := usecase.NewUserUsecaseImpl(userRepo, redisImpl)

	user := &model.User{
		ID:              "userID_1",
		FullName:        "ibanrama",
		Gender:          "undefined",
		Image:           ".png",
		Username:        "rama",
		Email:           "ibanrama29@gmail.com",
		Password:        "$2a$14$bUlfMzLFLKrILX3uNBkEieHMt/0boFfBvzpyDN4ai.UytAaTwmP1u",
		PhoneNumber:     sql.NullString{},
		EmailVerifiedAt: false,
		CreatedAt:       0,
		CreatedBy:       "userID_1",
		UpdatedAt:       0,
		UpdatedBy:       sql.NullString{},
		DeletedAt:       sql.NullInt64{},
		DeletedBy:       sql.NullString{},
	}

	t.Run("SUCCESS", func(t *testing.T) {
		_ = userRepo.OpenConn(context.TODO())
		defer userRepo.CloseConn()

		_ = userRepo.StartTx(context.TODO(), &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
			ReadOnly:  false,
		})
		userRepo.StartTxReturns(nil)
		defer func() {
			_ = userRepo.EndTx(nil)
			userRepo.EndTxReturns(nil)
		}()

		_, _ = userRepo.GetUserByEmail(context.TODO(), "ibanrama29@gmail.com")
		userRepo.GetUserByEmailReturns(user, nil)

		_ = userRepo.UpdateActivasiUser(context.TODO(), user)
		userRepo.UpdateActivasiUserReturns(nil)

		activasiResp, err := userUsecase.ActivasiAccount(context.TODO(), "ibanrama29@gmail.com")
		assert.NoError(t, err)
		assert.NotNil(t, activasiResp)
		assert.Equal(t, true, activasiResp.EmailVerifiedAt)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("ERROR_Invalid-email", func(t *testing.T) {
		_ = userRepo.OpenConn(context.TODO())
		defer userRepo.CloseConn()

		_ = userRepo.StartTx(context.TODO(), &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
			ReadOnly:  false,
		})
		userRepo.StartTxReturns(nil)
		defer func() {
			_ = userRepo.EndTx(nil)
			userRepo.EndTxReturns(nil)
		}()

		_, _ = userRepo.GetUserByEmail(context.TODO(), "nil")
		userRepo.GetUserByEmailReturns(nil, sql.ErrNoRows)

		activasiResp, err := userUsecase.ActivasiAccount(context.TODO(), "nil")
		assert.Error(t, err)
		assert.Nil(t, activasiResp)
		var errHTTP *model.ErrResponseHTTP
		assert.Equal(t, true, errors.As(err, &errHTTP))
		assert.Equal(t, 404, errHTTP.Code)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("ERROR_email-is-active", func(t *testing.T) {
		_ = userRepo.OpenConn(context.TODO())
		defer userRepo.CloseConn()

		_ = userRepo.StartTx(context.TODO(), &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
			ReadOnly:  false,
		})
		userRepo.StartTxReturns(nil)
		defer func() {
			_ = userRepo.EndTx(nil)
			userRepo.EndTxReturns(nil)
		}()

		user.EmailVerifiedAt = true
		_, _ = userRepo.GetUserByEmail(context.TODO(), "ibanrama29@gmail.com")
		userRepo.GetUserByEmailReturns(user, nil)

		activasiResp, err := userUsecase.ActivasiAccount(context.TODO(), "ibanrama29@gmail.com")
		assert.Error(t, err)
		assert.Nil(t, activasiResp)
		var errHTTP *model.ErrResponseHTTP
		assert.Equal(t, true, errors.As(err, &errHTTP))
		assert.Equal(t, 400, errHTTP.Code)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
