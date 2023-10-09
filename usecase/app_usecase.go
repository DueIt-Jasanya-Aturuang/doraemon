package usecase

import (
	"context"
)

type AppUsecase interface {
	CheckByID(ctx context.Context, appID string) error
}
