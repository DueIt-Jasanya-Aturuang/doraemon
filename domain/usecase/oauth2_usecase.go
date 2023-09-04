package usecase

import (
	"context"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/dto"
)

type Oauth2Usecase interface {
	GoogleClaimUser(ctx context.Context, req *dto.LoginGoogleReq) (*dto.LoginGoogleResp, error)
}
