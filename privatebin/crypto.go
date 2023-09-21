package privatebin

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/pbkdf2"
)

func deriveKey(key string, salt string, iterations int) ([]byte, error) {
	rawKey := base58.Decode(key)

	rawSalt, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return nil, err
	}

	return pbkdf2.Key(rawKey, rawSalt, iterations, 256/8, sha256.New), nil
}

func decryptContent(key []byte, iv string, authData any, content string) ([]byte, error) {
	rawIv, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, err
	}

	rawContent, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return nil, err
	}

	rawAuthData, err := json.Marshal(authData)
	if err != nil {
		return nil, err
	}

	block, _ := aes.NewCipher(key)
	aesGcm, _ := cipher.NewGCMWithNonceSize(block, len(rawIv))
	plainText, err := aesGcm.Open(nil, rawIv, rawContent, rawAuthData)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}
