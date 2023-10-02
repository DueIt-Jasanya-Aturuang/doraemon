package converter

import (
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain"
)

func GoogleClaimModelToResp(g *domain.Oauth2GoogleUser, exist bool) *domain.ResponseLoginWithGoogle {
	return &domain.ResponseLoginWithGoogle{
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
