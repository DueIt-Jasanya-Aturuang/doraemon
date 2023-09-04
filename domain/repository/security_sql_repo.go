package repository

import (
	"context"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/domain/model"
)

//counterfeiter:generate -o ./../mocks . SecuritySqlRepo
type SecuritySqlRepo interface {
	CreateToken(ctx context.Context, token *model.Token) (*model.Token, error)
	GetTokenByIDAndUserID(ctx context.Context, tokenID string, userID string) (*model.Token, error)
	UpdateToken(ctx context.Context, token *model.TokenUpdate) error
	DeleteToken(ctx context.Context, tokenID string, userID string) error
	DeleteAllTokenByUserID(ctx context.Context, userID string) error
	UnitOfWorkSqlRepo
}
