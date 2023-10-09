package repository

import (
	"context"
)

type SecurityRepository interface {
	Create(ctx context.Context, token *Security) error
	GetByAccessToken(ctx context.Context, token string) (*Security, error)
	Update(ctx context.Context, param UpdateSecurityParams) error
	Delete(ctx context.Context, id int, userID string) error
	DeleteAllByUserID(ctx context.Context, userID string) error
	UnitOfWorkRepository
}

type Security struct {
	ID           int
	UserID       string
	AppID        string
	RememberMe   bool
	AcceesToken  string
	RefreshToken string
}

type UpdateSecurityParams struct {
	ID           int
	RefreshToken string
	AccessToken  string
}
