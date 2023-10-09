package app_usecase

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
)

func (a *AppUsecaseImpl) CheckByID(ctx context.Context, appID string) error {
	exist, err := a.appRepo.CheckAppByID(ctx, appID)
	if err != nil {
		return err
	}

	if !exist {
		log.Info().Msgf("invalid app_repository id | data : %s", appID)
		return usecase.InvalidAppID
	}

	return nil
}
