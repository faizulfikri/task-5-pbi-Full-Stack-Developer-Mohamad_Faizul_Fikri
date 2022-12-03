package helperhash

import (
	"golang.org/x/crypto/bcrypt"
)

func GetHashPassword(pass string) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
	return bytes, err
}

func ComparePassword(hash, pass string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
}
