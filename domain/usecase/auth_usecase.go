package usecase

import (
	"context"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
)

type AuthUsecase interface {
	Login(ctx context.Context, req *dto.LoginReq) (*dto.UserResp, *dto.ProfileResp, error)
	Register(ctx context.Context, req *dto.RegisterReq) (*dto.UserResp, *dto.ProfileResp, error)
}
