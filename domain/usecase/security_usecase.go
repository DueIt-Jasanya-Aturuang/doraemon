package usecase

import (
	"context"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
)

type SecurityUsecase interface {
	JwtValidateAT(ctx context.Context, req *dto.JwtTokenReq, endpoint string) error
	JwtGenerateRTAT(ctx context.Context, req *dto.JwtTokenReq) (*dto.JwtTokenResp, error)
	ResetPassword(ctx context.Context, req *dto.ResetPasswordReq) error
	ForgottenPassword(ctx context.Context, req *dto.ForgottenPasswordReq) error
}
