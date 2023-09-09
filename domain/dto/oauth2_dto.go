package dto

type LoginGoogleReq struct {
	Token  string `json:"token"`
	Device string `json:"device"`
}

type LoginGoogleResp struct {
	ID            string
	Email         string
	VerifiedEmail bool
	Name          string
	GivenName     string
	FamilyName    string
	Image         string
	Locale        string
	ExistsUser    bool
}
