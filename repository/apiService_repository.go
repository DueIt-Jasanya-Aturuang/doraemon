package repository

type ApiServiceRepository interface {
	CreateProfileDueit(data []byte, appID string) (*ProfileDueit, error)
	GetProfileByUserIDDueit(userID string, appID string) (*ProfileDueit, error)
}

type ProfileDueit struct {
	ProfileID string  `json:"profile_id"`
	Quote     *string `json:"quote"`
	Profesi   *string `json:"profesi"`
}
