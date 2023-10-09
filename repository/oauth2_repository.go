package repository

type Oauth2Repository interface {
	GetGoogleToken(code string) (*Oauth2GoogleToken, error)
	GetGoogleUser(token *Oauth2GoogleToken) (*Oauth2GoogleUser, error)
}

type Oauth2GoogleUser struct {
	ID            string
	Email         string
	VerifiedEmail bool
	Name          string
	GivenName     string
	FamilyName    string
	Image         string
	Locale        string
}

type Oauth2GoogleToken struct {
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token"`
}
