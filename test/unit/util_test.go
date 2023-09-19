package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/util"
)

func TestRandomChar(t *testing.T) {
	rand, err := util.RandomChar(6)
	assert.NoError(t, err)
	t.Log(rand)
	assert.Equal(t, 6, len(rand))
}
