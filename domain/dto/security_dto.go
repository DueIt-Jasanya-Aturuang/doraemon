package dto

import (
	"time"
)

type JwtTokenReq struct {
	AppId         string
	Authorization string
	UserId        string
}

type JwtRegisteredTokenReq struct {
	AppId      string
	UserId     string
	RememberMe bool
}

type JwtTokenResp struct {
	RememberMe bool
	Token      string
	Exp        time.Duration
}
