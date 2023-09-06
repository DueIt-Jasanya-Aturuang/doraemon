package usecase

import (
	"context"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
)

type AppUsecase interface {
	CheckAppByID(ctx context.Context, req *dto.AppReq) error
}
