package domain

import (
	"context"
)

//counterfeiter:generate -o ./../mocks . SecuritySqlRepo
type SecurityRepository interface {
	Create(ctx context.Context, token *Token) error
	GetByAccessToken(ctx context.Context, token string) (*Token, error)
	Update(ctx context.Context, id int, refreshToken string, accessToken string) error
	Delete(ctx context.Context, id int, userID string) error
	DeleteAllByUserID(ctx context.Context, userID string) error
	UnitOfWorkRepository
}

type SecurityUsecase interface {
	JwtValidation(ctx context.Context, req *RequestJwtToken) (bool, error)
	JwtGenerate(ctx context.Context, req *RequestJwtToken) (*ResponseJwtToken, error)
	Logout(ctx context.Context, req *RequestLogout) (err error)
}

type Token struct {
	ID           int
	UserID       string
	AppID        string
	RememberMe   bool
	AcceesToken  string
	RefreshToken string
}

type RequestJwtToken struct {
	AppId          string
	Authorization  string
	UserId         string
	ActivasiHeader bool
}

type ResponseJwtToken struct {
	Token string `json:"token"`
}

type RequestLogout struct {
	Token  string
	UserID string
}
