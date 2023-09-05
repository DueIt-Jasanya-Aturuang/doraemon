package usecase

import (
	"context"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
)

type SecurityUsecase interface {
	JwtValidateAT(ctx context.Context, req *dto.JwtTokenReq, endpoint string) (string, bool, error)
	JwtGenerateRTAT(ctx context.Context, req *dto.JwtTokenReq) (*dto.JwtTokenResp, error)
	JwtRegistredRTAT(ctx context.Context, req *dto.JwtRegisteredTokenReq) (*dto.JwtTokenResp, error)
}
