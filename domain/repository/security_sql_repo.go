package repository

import (
	"context"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
)

//counterfeiter:generate -o ./../mocks . SecuritySqlRepo
type SecuritySqlRepo interface {
	CreateToken(ctx context.Context, token *model.Token) error
	GetTokenByAT(ctx context.Context, token string) (*model.Token, error)
	UpdateToken(ctx context.Context, id int, refreshToken string, accessToken string) error
	DeleteToken(ctx context.Context, id int, userID string) error
	DeleteAllTokenByUserID(ctx context.Context, userID string) error
	UnitOfWorkSqlRepo
}
