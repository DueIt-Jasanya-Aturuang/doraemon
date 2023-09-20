package domain

import (
	"context"
)

//counterfeiter:generate -o ./../mocks . AppSqlRepo
type AppRepository interface {
	CheckAppByID(ctx context.Context, id string) (bool, error)
	UnitOfWorkRepository
}

type AppUsecase interface {
	CheckByID(ctx context.Context, req *RequestCheckApp) error
}

type RequestCheckApp struct {
	AppID string
}
