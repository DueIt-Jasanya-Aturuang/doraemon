package integration

import (
	"database/sql"
	"testing"
)

var db *sql.DB

func TestInit(t *testing.T) {
	t.Run("AccountApiRepo_CREATE_PROFILE", CreateProfile)
}
