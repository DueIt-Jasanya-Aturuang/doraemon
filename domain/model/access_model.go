package model

import "database/sql"

type Access struct {
	ID             int8
	AppId          string
	UserId         string
	RoleId         int8
	AccessEndpoint string
	CreatedAt      int64
	CreatedBy      string
	UpdatedAt      int64
	UpdatedBy      sql.NullString
	DeletedAt      sql.NullInt64
	DeletedBy      sql.NullString
}
