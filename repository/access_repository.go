package repository

import (
	"context"
)

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
