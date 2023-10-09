package usecase

type ApiServiceUsecase interface {
	GetProfileByUserID(userID string, appID string) (*ResponseProfileDueit, error)
}

type ResponseProfileDueit struct {
	ProfileID string  `json:"profile_id"`
	Quote     *string `json:"quote"`
	Profesi   *string `json:"profesi"`
}

type RequestCreateProfile struct {
	UserID string `json:"user_id"`
}
