package encryption

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecryptStringCFB(t *testing.T) {
	key := "rp1LZ7WNMXaq9KhK2LIkMbAA8qbMkFQb"
	text := "6fe9168f11b0ea3cf529e3363d7e74f09f95a84b"
	res, err := DecryptStringCFB(text, key)
	assert.NoError(t, err)
	t.Log(res)
}

func TestEncrypStringCFB(t *testing.T) {
	key := "rp1LZ7WNMXaq9KhK2LIkMbAA8qbMkFQb"
	res, err := EncrypStringCFB("rama", key)
	assert.NoError(t, err)
	t.Log(res)
}
func TestDecryptStringCBC(t *testing.T) {
	key := "rp1LZ7WNMXaq9KhK2LIkMbAA8qbMkFQb"
	iv := "my16digitIvKey12"
	text := "sxy5+KAARCFYkU16AQPWyw=="
	res, err := DecryptStringCBC(text, key, iv)
	assert.NoError(t, err)
	t.Log(res)
}

func TestEncrypStringCBC(t *testing.T) {
	key := "rp1LZ7WNMXaq9KhK2LIkMbAA8qbMkFQb"
	iv := "my16digitIvKey12"
	res, err := EncrypStringCBC("rama", key, iv)
	assert.NoError(t, err)
	t.Log(res)
}
