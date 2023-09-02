package dto

type OTPGenerateReq struct {
	Email string `json:"email" validate:"required,email,min=3,max=55"`
}

type OTPValidationReq struct {
	Email string `json:"email" validate:"required,email,min=3,max=55"`
	OTP   string `json:"otp" validate:"required,min=6,max=6"`
}
