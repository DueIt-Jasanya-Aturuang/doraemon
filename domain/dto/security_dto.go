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

type JwtRegisteredTokenReq struct {
	AppId         string
	Authorization string
	UserId        string
	Activasi      bool
	RememberMe    bool
}

type JwtTokenResp struct {
	RememberMe bool
	Token      string
	Exp        time.Duration
}
