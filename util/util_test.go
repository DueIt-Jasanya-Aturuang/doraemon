package util

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomChar(t *testing.T) {
	for i := 1; i <= 1000; i++ {
		if i == 1000 {
			Charset += strconv.Itoa(i)
		} else {
			Charset += strconv.Itoa(i) + ""
		}
	}

	user := "userGoogle"
	randomChar, err := RandomChar(6)
	assert.NoError(t, err)
	t.Logf(user + randomChar)
}
