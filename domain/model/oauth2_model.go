package model

type GoogleOauth2Token struct {
	AccessToken string
	IDToken     string
}

type GoogleOauth2User struct {
	ID            string
	Email         string
	VerifiedEmail bool
	Name          string
	GivenName     string
	FamilyName    string
	Image         string
	Locale        string
}
