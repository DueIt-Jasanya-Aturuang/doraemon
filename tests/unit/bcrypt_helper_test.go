package unit

import (
	"github.com/DueIt-Jasanya-Aturuang/doraemon/internal/helper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBcryptHelper(t *testing.T) {
	password := "rama"
	hash := ""

	t.Run("SUCCESS_hash", func(t *testing.T) {
		pw, err := helper.BcryptPasswordHash(password)
		assert.NoError(t, err)
		t.Log(pw)
		hash = pw
	})

	t.Run("SUCCESS_compare", func(t *testing.T) {
		same := helper.BcryptPasswordCompare(password, hash)
		assert.Equal(t, true, same)
	})

	t.Run("ERROR_compare", func(t *testing.T) {
		hash = "any passsword"
		same := helper.BcryptPasswordCompare(password, hash)
		assert.Equal(t, false, same)
	})

}
