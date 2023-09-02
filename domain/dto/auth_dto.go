package dto

type RegisterReq struct {
	FullName        string `json:"full_name" validate:"required,min=3,max=32"`
	Username        string `json:"username" validate:"required,min=3,max=22"`
	Email           string `json:"email" validate:"required,email,min=3,max=55"`
	Password        string `json:"password" validate:"required,min=6,max=55"`
	RePassword      string `json:"re_password" validate:"required,min=6,max=55"`
	EmailVerifiedAt bool   `swaggerignore:"true"`
	Oauth           bool   `swaggerignore:"true"`
}

type LoginReq struct {
	EmailOrUsername string `json:"email_or_username" validate:"required,min=3,max=55"`
	Password        string `json:"password" validate:"required,min=6,max=55"`
	RememberMe      bool   `json:"remember_me" validate:"required"`
	Oauth           bool   `swaggerignore:"true"`
}

type LogoutReq struct {
	AppId         string
	Authorization string
	UserId        string
	Activasi      bool
}
