package usecase

import (
	"context"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
)

type UserUsecase interface {
	ResetPassword(ctx context.Context, req *dto.ResetPasswordReq) error
	ForgottenPassword(ctx context.Context, req *dto.ForgottenPasswordReq) (string, error)
	ResetForgottenPassword(ctx context.Context, req *dto.ResetForgottenPasswordReq) error
	ActivasiAccount(c context.Context, email string) (resp *dto.ActivasiAccountResp, err error)
}
