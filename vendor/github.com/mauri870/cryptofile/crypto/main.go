package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

var (
	ErrorCypherTooShort   = errors.New("Ciphertext too short")
	ErrorKeyInvalidLength = errors.New("The key length must have 32 bytes for AES-256 or 16 bytes for AES-128")
)

func Encrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, aes.BlockSize+len(string(text)))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], text)

	return []byte(base64.StdEncoding.EncodeToString(ciphertext)), nil
}

func Decrypt(key, ciphertext []byte) ([]byte, error) {
	text, _ := base64.StdEncoding.DecodeString(string(ciphertext))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, ErrorCypherTooShort
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)

	return text, nil
}
