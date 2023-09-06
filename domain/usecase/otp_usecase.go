package usecase

import (
	"context"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
)

//counterfeiter:generate -o ./../mocks . OTPUsecase
type OTPUsecase interface {
	OTPGenerate(ctx context.Context, req *dto.OTPGenerateReq) error
	OTPValidation(ctx context.Context, req *dto.OTPValidationReq) error
}
