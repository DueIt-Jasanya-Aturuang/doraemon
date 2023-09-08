package dto

type OTPGenerateReq struct {
	Email  string `json:"email"`
	Type   string // Type set in header
	UserID string // UserID set in header
}

type OTPValidationReq struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
	Type  string // Type set in handler
}
