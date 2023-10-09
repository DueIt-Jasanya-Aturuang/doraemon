package app_usecase

import (
	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/usecase"
)

type AppUsecaseImpl struct {
	appRepo repository.AppRepository
}

func NewAppUsecaseImpl(
	appRepo repository.AppRepository,
) usecase.AppUsecase {
	return &AppUsecaseImpl{
		appRepo: appRepo,
	}
}
