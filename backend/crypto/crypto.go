package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

var (
	passphrase []byte
)

func InitCrypto() {
	secretPath := "/config/aes.json"

	if altp, ok := os.LookupEnv("AES_SECRETS"); ok {
		secretPath = altp
	}

	f, err := ioutil.ReadFile(secretPath)

	if err != nil {
		panic(err)
	}

	if len(f) != 32 {
		panic(fmt.Errorf("AES Key should be 32 bytes of length"))
	}

	passphrase = f
}

func Key() []byte {
	return passphrase
}

func Encrypt(data []byte) []byte {
	block, _ := aes.NewCipher(passphrase)

	gcm, err := cipher.NewGCM(block)

	if err != nil {
		panic(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func Decrypt(data []byte) []byte {
	block, err := aes.NewCipher(passphrase)

	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(block)

	if err != nil {
		panic(err)
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)

	if err != nil {
		panic(err)
	}

	return plaintext
}
