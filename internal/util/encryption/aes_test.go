package encryption

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	iv := "auTHxeaZSrsxfIZI"
	text := "sxy5+KAARCFYkU16AQPWyw=="
	res, err := DecryptStringCBC(text, key, iv)
	assert.NoError(t, err)
	t.Log(res)
}

func TestEncrypStringCBC(t *testing.T) {
	key := "rp1LZ7WNMXaq9KhK2LIkMbAA8qbMkFQb"
	iv := "mydigit15digit11"

	// req := map[string]string{
	// 	"access_token": "this is access token",
	// 	"id_token":     "this is access token",
	// }
	//
	// reqJson, _ := json.Marshal(req)

	res, err := EncrypStringCBC("rama", key, iv)
	assert.NoError(t, err)
	t.Log(res)
}
