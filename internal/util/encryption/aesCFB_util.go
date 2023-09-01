package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
)

func DecryptStringCFB(text string, key string) (string, error) {
	if len(key) > 32 {
		log.Warn().Msg("key must be 32 character")
		return "", fmt.Errorf("%s", "invalid key encryptions")
	}

	cipherText, err := hex.DecodeString(text)
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

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}

func EncrypStringCFB(text string, key string) (string, error) {
	if len(key) > 32 {
		log.Warn().Msg("key must be 32 character")
		return "", fmt.Errorf("%s", "invalid key encryptions")
	}

	plainText := []byte(text)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Err(err).Msg("failed to generate block aes")
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]

	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		log.Err(err).Msg("failed to generate initalization vector")
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return fmt.Sprintf("%x", cipherText), nil
}
