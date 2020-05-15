package models

import (
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 13

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func TokenHash(id string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(id), bcryptCost)
	if err != nil {
		return ""
	}
	hasher := sha256.New()
	hasher.Write(hash)
	return hex.EncodeToString(hasher.Sum(nil))
}
