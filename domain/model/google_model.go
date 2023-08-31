package model

type GoogleOauthToken struct {
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token"`
}

type GoogleOauthUser struct {
	ID            string
	Email         string
	VerifiedEmail bool
	Name          string
	GivenName     string
	FamilyName    string
	Image         string
	Locale        string
}
