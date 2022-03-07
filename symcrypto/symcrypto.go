package symcrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"golang.org/x/crypto/sha3"
)

func Encrypt(plaintext string, passphrase string) []byte {
	aesgcm := newGCM(passphrase)

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	ciphertext := []byte{}
	ciphertext = append(ciphertext, nonce...)
	ciphertext = append(ciphertext, aesgcm.Seal(nil, nonce, []byte(plaintext), nil)...)

	return ciphertext
}

func Decrypt(ciphertext []byte, passphrase string) string {
	aesgcm := newGCM(passphrase)

	nonceSize := aesgcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return string(plaintext)
}

func newGCM(passphrase string) cipher.AEAD {
	sum := sha3.Sum256([]byte(passphrase))
	key := sum[:]

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	return aesgcm
}
