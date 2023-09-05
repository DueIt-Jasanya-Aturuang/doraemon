package dto

type OTPGenerateReq struct {
	Email string `json:"email" validate:"required,email,min=3,max=55"`
	Type  string // Type set in header
}

type OTPValidationReq struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
	Type  string // Type set in header
}
