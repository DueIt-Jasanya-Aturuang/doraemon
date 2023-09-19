package domain

type MicroServiceRepository interface {
	CreateProfile(data []byte) (*Profile, error)
	GetProfileByUserID(userID string) (*Profile, error)
}

type Profile struct {
	ProfileID string `json:"profile_id"`
	Quote     string `json:"quote"`
	Profesi   string `json:"profesi"`
}

type RequestCreateProfile struct {
	UserID string `json:"user_id"`
}

type ResponseProfile struct {
	ProfileID string `json:"profile_id"`
	Quote     string `json:"quote"`
	Profesi   string `json:"profesi"`
}
