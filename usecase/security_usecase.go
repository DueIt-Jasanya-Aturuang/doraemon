package usecase

import (
	"context"
)

type SecurityUsecase interface {
	ValidationJWT(ctx context.Context, req *RequestValidationJWT) (bool, error)
	ReGenerateJWT(ctx context.Context, req *RequestValidationJWT) (*ResponseJWT, error)
	GenerateJWT(ctx context.Context, req *RequestGenerateJWT) (*ResponseJWT, error)
	Logout(ctx context.Context, req *RequestLogout) (err error)
}

type RequestValidationJWT struct {
	AppID          string
	Authorization  string
	UserID         string
	ActivasiHeader bool
}

type RequestGenerateJWT struct {
	AppID      string
	UserID     string
	RememberMe bool
}

type ResponseJWT struct {
	Token string
}

type RequestLogout struct {
	Token  string
	UserID string
}
