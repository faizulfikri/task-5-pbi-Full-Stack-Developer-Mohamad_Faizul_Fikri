package otentifikasi

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type ClaimJwt struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

func GenerateJWT(email string, username string) (tokenStr string, err error) {
	mySigningKey := []byte("AllForOne")
	expiredTime := time.Now().Add(1 * time.Hour)

	claims := &ClaimJwt{
		Email:    email,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err = token.SignedString(mySigningKey)

	return
}
