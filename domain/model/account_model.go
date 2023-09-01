package model

type Profile struct {
	Code      int
	ProfileID string `json:"profile_id"`
	Quote     string `json:"quote"`
	Profesi   string `json:"profesi"`
}
