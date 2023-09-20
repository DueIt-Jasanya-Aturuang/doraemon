package domain

import (
	"context"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate -o ./../mocks . AccessSqlRepo
type AccessRepository interface {
	Create(ctx context.Context, access *Access) error
	GetByUserIDAndAppID(ctx context.Context, userID string, appID string) (*Access, error)
	UnitOfWorkRepository
}

type Access struct {
	ID             int8
	AppId          string
	UserId         string
	RoleId         int8
	AccessEndpoint string
	AuditInfo
}
