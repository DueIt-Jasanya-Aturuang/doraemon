package unit

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/mocks"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
)

func TestAppUsecaseCheckAppByID(t *testing.T) {
	appRepo := &mocks.FakeAppSqlRepo{}

	appUsecase := usecase.NewAppUsecaseImpl(appRepo)

	t.Run("SUCCESS", func(t *testing.T) {
		err := appRepo.OpenConn(context.TODO())
		assert.NoError(t, err)
		defer appRepo.CloseConn()

		_, err = appRepo.CheckAppByID(context.TODO(), "123")
		assert.NoError(t, err)
		appRepo.CheckAppByIDReturns(true, nil)

		req := &dto.AppReq{
			AppID: "123",
		}
		err = appUsecase.CheckAppByID(context.TODO(), req)
		assert.NoError(t, err)
		assert.Equal(t, "123", req.AppID)
	})

	t.Run("ERROR", func(t *testing.T) {
		err := appRepo.OpenConn(context.TODO())
		assert.NoError(t, err)
		defer appRepo.CloseConn()

		_, err = appRepo.CheckAppByID(context.TODO(), "1233")
		assert.NoError(t, err)
		appRepo.CheckAppByIDReturns(false, nil)

		req := &dto.AppReq{
			AppID: "123",
		}
		err = appUsecase.CheckAppByID(context.TODO(), req)
		assert.Error(t, err)
		var errorHTTP *model.ErrResponseHTTP
		assert.Equal(t, true, errors.As(err, &errorHTTP))
		t.Log(errorHTTP.Code)
	})
}
