package usecase

import (
	"context"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
)

//counterfeiter:generate -o ./../mocks . SecurityUsecase
type SecurityUsecase interface {
	JwtValidateAT(ctx context.Context, req *dto.JwtTokenReq, activasiHeader string) (bool, error)
	JwtGenerateRTAT(ctx context.Context, req *dto.JwtTokenReq) (*dto.JwtTokenResp, error)
	Logout(ctx context.Context, req *dto.LogoutReq) (err error)
}
