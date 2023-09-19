package domain

import (
	"context"
)

type OTPUsecase interface {
	Generate(ctx context.Context, req *RequestGenerateOTP) error
	Validation(ctx context.Context, req *RequestValidationOTP) error
}

type RequestGenerateOTP struct {
	Email  string `json:"email"`
	Type   string // Type set in header
	UserID string // UserID set in header
}

type RequestValidationOTP struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
	Type  string // Type set in handler
}
