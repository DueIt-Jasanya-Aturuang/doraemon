package usecase

import (
	"context"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
)

//counterfeiter:generate -o ./../mocks . SecurityUsecase
type SecurityUsecase interface {
	JwtValidateAT(ctx context.Context, req *dto.JwtTokenReq, endpoint string) (bool, error)
	JwtGenerateRTAT(ctx context.Context, req *dto.JwtTokenReq) (*dto.JwtTokenResp, error)
	JwtRegistredRTAT(ctx context.Context, req *dto.JwtRegisteredTokenReq) (*dto.JwtTokenResp, error)
	Logout(ctx context.Context, req *dto.LogoutReq) (err error)
}
