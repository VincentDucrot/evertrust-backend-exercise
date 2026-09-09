package service

import (
	"crypto/rand"
	"crypto/sha256"
)

func NewToken() (plaintext string, hash [32]byte) {
	plaintext = rand.Text()
	hash = sha256.Sum256([]byte(plaintext))
	return plaintext, hash
}
