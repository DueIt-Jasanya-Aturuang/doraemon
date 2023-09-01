package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"

	"github.com/rs/zerolog/log"
)

func DecryptStringCBC(text string, key string, iv string) (string, error) {
	if len(key) > 32 {
		log.Warn().Msg("key must be 32 character")
		return "", fmt.Errorf("%s", "invalid key encryptions")
	}

	cipherText, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		log.Err(err).Msg("failed to decord string")
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Err(err).Msg("failed to generate block aes")
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		return "", fmt.Errorf("%s", "cipher text is too short")
	}

	mode := cipher.NewCBCDecrypter(block, []byte(iv))
	mode.CryptBlocks(cipherText, cipherText)

	cipherText = unpaddingPKCS7(cipherText)

	return string(cipherText), nil
}

func paddingPKCS7(plainText []byte) []byte {
	padding := aes.BlockSize - len(plainText)%aes.BlockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(plainText, padText...)
}

func unpaddingPKCS7(plaintext []byte) []byte {
	length := len(plaintext)
	unpadding := int(plaintext[length-1])
	return plaintext[:(length - unpadding)]
}
func EncrypStringCBC(text string, key string, iv string) (string, error) {
	if len(key) > 32 {
		log.Warn().Msg("key must be 32 character")
		return "", fmt.Errorf("%s", "invalid key encryptions")
	}

	plainText := []byte(text)
	plainText = paddingPKCS7(plainText)

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Err(err).Msg("failed to generate block aes")
		return "", err
	}

	cipherText := make([]byte, len(plainText))

	mode := cipher.NewCBCEncrypter(block, []byte(iv))
	mode.CryptBlocks(cipherText, plainText)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}
