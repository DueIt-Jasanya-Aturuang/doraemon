package helper

import (
	"database/sql"
)

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
