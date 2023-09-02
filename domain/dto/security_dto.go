package dto

import (
	"time"
)

type JwtTokenReq struct {
	AppId         string
	Authorization string
	UserId        string
	Activasi      bool
}

type JwtTokenResp struct {
	RememberMe bool
	Token      string
	Exp        time.Duration
}

type ResetPasswordReq struct {
	OldPassword string `json:"old_password" validate:"required,min=6,max=55"`
	Password    string `json:"password" validate:"required,min=6,max=55"`
	RePassword  string `json:"re_password" validate:"required,min=6,max=55"`
}

type ForgottenPasswordReq struct {
	OldPassword string `json:"old_password" validate:"required,min=6,max=55"`
	Password    string `json:"password" validate:"required,min=6,max=55"`
	RePassword  string `json:"re_password" validate:"required,min=6,max=55"`
}
