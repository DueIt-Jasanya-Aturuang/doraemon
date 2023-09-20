package domain

import (
	"context"
)

type Oauth2Repository interface {
	GetGoogleToken(code string) (*Oauth2GoogleToken, error)
	GetGoogleUser(token *Oauth2GoogleToken) (*Oauth2GoogleUser, error)
}

type Oauth2Usecase interface {
	GoogleClaimUser(ctx context.Context, req *RequestLoginWithGoogle) (*ResponseLoginWithGoogle, error)
}

type Oauth2GoogleToken struct {
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token"`
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

type ResponseLoginWithGoogle struct {
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

type RequestLoginWithGoogle struct {
	Token  string `json:"token"`
	Device string `json:"device"`
}
