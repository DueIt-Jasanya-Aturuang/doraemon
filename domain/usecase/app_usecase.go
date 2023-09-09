package usecase

import (
	"context"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate -o ./../mocks . AppUsecase
type AppUsecase interface {
	CheckAppByID(ctx context.Context, req *dto.AppReq) error
}
