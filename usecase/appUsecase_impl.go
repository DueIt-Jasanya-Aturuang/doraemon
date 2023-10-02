package usecase

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain"
)

type AppUsecaseImpl struct {
	appRepo domain.AppRepository
}

func NewAppUsecaseImpl(
	appRepo domain.AppRepository,
) domain.AppUsecase {
	return &AppUsecaseImpl{
		appRepo: appRepo,
	}
}

func (a *AppUsecaseImpl) CheckByID(ctx context.Context, req *domain.RequestCheckApp) error {
	if err := a.appRepo.OpenConn(ctx); err != nil {
		return err
	}
	defer a.appRepo.CloseConn()

	exist, err := a.appRepo.CheckAppByID(ctx, req.AppID)
	if err != nil {
		return err
	}

	if !exist {
		log.Info().Msgf("invalid app id | data : %s", req.AppID)
		return InvalidAppID
	}

	return nil
}
