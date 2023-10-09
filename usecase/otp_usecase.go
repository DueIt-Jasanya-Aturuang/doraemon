package usecase

import (
	"context"
)

type OTPUsecase interface {
	Generate(ctx context.Context, req *RequestGenerateOTP) error
	Validation(ctx context.Context, req *RequestValidationOTP) error
}

type RequestGenerateOTP struct {
	Email  string
	Type   string
	UserID string
}

type RequestValidationOTP struct {
	Email string
	OTP   string
	Type  string
}
