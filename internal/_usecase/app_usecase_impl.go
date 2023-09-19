package _usecase

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/usecase"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/util/error"
)

type AppUsecaseImpl struct {
	appRepo repository.AppSqlRepo
}

func NewAppUsecaseImpl(
	appRepo repository.AppSqlRepo,
) usecase.AppUsecase {
	return &AppUsecaseImpl{
		appRepo: appRepo,
	}
}

func (a *AppUsecaseImpl) CheckAppByID(ctx context.Context, req *dto.AppReq) (err error) {
	// OpenConn open connection dari app repo
	// defer untuk close connection
	err = a.appRepo.OpenConn(ctx)
	if err != nil {
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}
	defer a.appRepo.CloseConn()

	// CheckAppByID apakah app id valid atau gak
	// jika tidak maka akan return 403
	exists, err := a.appRepo.CheckAppByID(ctx, req.AppID)
	if err != nil {
		return _error.ErrStringDefault(http.StatusInternalServerError)
	}
	if !exists {
		log.Info().Msgf("app id is not registered: %s", req.AppID)
		return _error.ErrStringDefault(http.StatusForbidden)
	}

	return nil
}
