package repository

import (
	"database/sql"
)

type AuditInfo struct {
	CreatedAt int64
	CreatedBy string
	UpdatedAt int64
	UpdatedBy sql.NullString
	DeletedAt sql.NullInt64
	DeletedBy sql.NullString
}

func NewNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}

	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func GetNullString(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}

	return nil
}
