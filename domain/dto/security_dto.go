package dto

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
	Token string `json:"token"`
}
