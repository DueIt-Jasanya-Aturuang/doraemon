package usecase

import (
	"context"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
)

type UserUsecase interface {
	ResetPassword(ctx context.Context, req *dto.ResetPasswordReq) error
	ForgottenPassword(ctx context.Context, req *dto.ForgottenPasswordReq) error
}
