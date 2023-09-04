package conv

import (
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
)

func GoogleClaimModelToResp(g *model.GoogleOauth2User, exist bool) *dto.LoginGoogleResp {
	return &dto.LoginGoogleResp{
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
