package domain

import (
	"context"
)

//counterfeiter:generate -o ./../mocks . AuthUsecase
type AuthUsecase interface {
	Login(ctx context.Context, req *RequestLogin) (*ResponseAuth, error)
	Register(ctx context.Context, req *RequestRegister) (*ResponseAuth, error)
}

type ResponseAuth struct {
	ResponseUser     `json:"user"`
	ResponseProfile  `json:"profile"`
	ResponseJwtToken `json:"token"`
}
type RequestRegister struct {
	FullName        string `json:"full_name"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	RePassword      string `json:"re_password"`
	EmailVerifiedAt bool   // EmailVerifiedAt set in handler
	AppID           string
	Role            int8 // Role set in handler
}

type RequestLogin struct {
	EmailOrUsername string `json:"email_or_username"`
	Password        string `json:"password"`
	RememberMe      bool   `json:"remember_me"`
	Oauth2          bool   // Oauth2 set in handler
	AppID           string // AppID set in handlr
}
