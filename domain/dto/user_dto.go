package dto

type UserResp struct {
	ID              string `json:"id"`
	FullName        string `json:"full_name"`
	Gender          string `json:"gender"`
	Image           string `json:"image"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	EmailFormat     string `json:"email_format"`
	PhoneNumber     string `json:"phone_number"`
	EmailVerifiedAt bool   `json:"activited"`
}

type ResetPasswordReq struct {
	OldPassword string `json:"old_password"`
	Password    string `json:"password"`
	RePassword  string `json:"re_password"`
	UserID      string // UserID get in header
}

type ForgottenPasswordReq struct {
	Email string // Email get in query param
}

type ResetForgottenPasswordReq struct {
	Email      string // Email get in query param
	Token      string // Token get in query param
	Password   string `json:"password"`
	RePassword string `json:"re_password"`
}

type ActivasiAccountResp struct {
	EmailVerifiedAt bool `json:"activited"`
}
