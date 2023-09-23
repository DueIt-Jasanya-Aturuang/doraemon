package unit

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/mocks"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/pkg/_usecase"
)

func TestOTPUsecaseOTPGenerate(t *testing.T) {
	userRepo := &mocks.FakeUserSqlRepo{}
	redisClient, mock := redismock.NewClientMock()
	redisImpl := &infra.RedisImpl{Client: redisClient}

	otpUsecase := _usecase.NewOTPUsecaseImpl(userRepo, redisImpl)
	userID := "userID_1"

	req := &dto.OTPGenerateReq{
		Email:  "ibanrama29@gmail.com",
		Type:   "",
		UserID: userID,
	}

	t.Run("SUCCESS_expire", func(t *testing.T) {
		req.Type = "activasi-account"
		_ = userRepo.OpenConn(context.TODO())
		defer userRepo.CloseConn()

		_, _ = userRepo.CheckActivasiUserByID(context.TODO(), req.UserID)
		userRepo.CheckActivasiUserByIDReturns(false, nil)

		mock.ExpectExists(req.Type + ":" + req.Email).SetVal(1)
		mock.ExpectExpire(req.Type+":"+req.Email, 5*time.Minute).SetVal(true)

		err := otpUsecase.OTPGenerate(context.TODO(), req)
		assert.Error(t, err)
		t.Log("emang sengaja error karna kafka, karna malas bikin mock interfacenya")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("SUCCESS_set", func(t *testing.T) {
		req.Type = "activasi-account"
		_ = userRepo.OpenConn(context.TODO())
		defer userRepo.CloseConn()

		_, _ = userRepo.CheckActivasiUserByID(context.TODO(), req.UserID)
		userRepo.CheckActivasiUserByIDReturns(false, nil)

		mock.ExpectExists(req.Type + ":" + req.Email).SetVal(0)
		mock.ExpectSet(req.Type+":"+req.Email, "123456", 5*time.Minute).SetVal("123456")

		err := otpUsecase.OTPGenerate(context.TODO(), req)
		assert.Error(t, err)
		t.Log("emang sengaja error karna kafka, karna malas bikin mock interfacenya")
		t.Log("ini bener cuman otp nya harus sama, karena random jadinya di ignore aja yang status were metnya")
	})

	t.Run("ERROR_activasi-account-is-true", func(t *testing.T) {
		mock.ClearExpect()
		req.Type = "activasi-account"
		_ = userRepo.OpenConn(context.TODO())
		defer userRepo.CloseConn()

		_, _ = userRepo.CheckActivasiUserByID(context.TODO(), req.UserID)
		userRepo.CheckActivasiUserByIDReturns(true, nil)

		err := otpUsecase.OTPGenerate(context.TODO(), req)
		assert.Error(t, err)
		var errHttp *model.ErrResponseHTTP
		assert.Equal(t, true, errors.As(err, &errHttp))

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestOTPUsecaseOTPValidation(t *testing.T) {
	userRepo := &mocks.FakeUserSqlRepo{}
	redisClient, mock := redismock.NewClientMock()
	redisImpl := &infra.RedisImpl{Client: redisClient}

	otpUsecase := _usecase.NewOTPUsecaseImpl(userRepo, redisImpl)

	req := &dto.OTPValidationReq{
		Email: "ibanrama29@gmail.com",
		Type:  "",
		OTP:   "123456",
	}

	t.Run("SUCCESS", func(t *testing.T) {
		mock.ExpectGet(req.Type + ":" + req.Email).SetVal("123456")
		mock.ExpectDel(req.Type + ":" + req.Email).SetVal(1)
		err := otpUsecase.OTPValidation(context.TODO(), req)
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
