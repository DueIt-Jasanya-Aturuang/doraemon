package repository

import (
	"context"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
)

type AccountApiRepo interface {
	CreateProfile(ctx context.Context, data []byte) (*model.Profile, error)
}
