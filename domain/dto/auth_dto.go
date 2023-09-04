package dto

type RegisterReq struct {
	FullName        string `json:"full_name"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	RePassword      string `json:"re_password"`
	EmailVerifiedAt bool   // EmailVerifiedAt set in handler
	AppId           string // AppId get in header
	Role            int8   // Role set in handler
}

type LoginReq struct {
	EmailOrUsername string `json:"email_or_username"`
	Password        string `json:"password"`
	RememberMe      bool   `json:"remember_me"`
	AppId           string // AppId get in header
	Oauth2          bool   // Oauth2 helper for login
}

type LogoutReq struct {
	Token string
}
