package dto

type RegisterReq struct {
	FullName        string `json:"full_name"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	RePassword      string `json:"re_password"`
	EmailVerifiedAt bool   // EmailVerifiedAt set in handler
	AppID           string
	Role            int8 // Role set in handler
}

type LoginReq struct {
	EmailOrUsername string `json:"email_or_username"`
	Password        string `json:"password"`
	RememberMe      bool   `json:"remember_me"`
	Oauth2          bool   // Oauth2 set in handler
	AppID           string // AppID set in handlr
}

type LogoutReq struct {
	Token  string
	UserID string
}

type ProfileReq struct {
	UserID string `json:"user_id"`
}

type ProfileResp struct {
	ProfileID string `json:"profile_id"`
	Quote     string `json:"quote"`
	Profesi   string `json:"profesi"`
}
