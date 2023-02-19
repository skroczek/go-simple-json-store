// Package helper
// @see https://www.thepolyglotdeveloper.com/2018/02/encrypt-decrypt-data-golang-application-crypto-packages/
package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
	"os"
)

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Encrypt(data []byte) []byte {
	passphrase, ok := os.LookupEnv("ACME_RESTFUL_PASSPHRASE")
	if !ok {
		panic("No passphrase set. Please set the ACME_RESTFUL_PASSPHRASE environment variable.")
	}
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	return gcm.Seal(nonce, nonce, data, nil)
}

func Decrypt(data []byte, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	passphrase, ok := os.LookupEnv("ACME_RESTFUL_PASSPHRASE")
	if !ok {
		panic("No passphrase set. Please set the ACME_RESTFUL_PASSPHRASE environment variable.")
	}
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
