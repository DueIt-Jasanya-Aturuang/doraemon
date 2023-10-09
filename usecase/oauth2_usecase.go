package usecase

import (
	"context"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/repository"
)

type Oauth2Usecase interface {
	GoogleClaimUser(ctx context.Context, req *RequestLoginWithGoogle) (*ResponseLoginWithGoogle, error)
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
	Token  string
	Device string
}

func GoogleClaimModelToResp(g *repository.Oauth2GoogleUser, exist bool) *ResponseLoginWithGoogle {
	return &ResponseLoginWithGoogle{
		ID:            g.ID,
		Email:         g.Email,
		VerifiedEmail: g.VerifiedEmail,
		Name:          g.Name,
		GivenName:     g.GivenName,
		FamilyName:    g.FamilyName,
		Image:         g.Image,
		Locale:        g.Locale,
		ExistsUser:    exist,
	}
}
