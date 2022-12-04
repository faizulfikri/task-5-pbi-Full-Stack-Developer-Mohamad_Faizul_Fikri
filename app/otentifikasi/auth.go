package otentifikasi

import (
	"errors"
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

func ValidateToken(tokenSign string) (err error) {
	mySigningKey := []byte("AllForOne")
	token, err := jwt.ParseWithClaims(
		tokenSign,
		&ClaimJwt{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(mySigningKey), nil
		},
	)

	if err != nil {
		return
	}
	claims, ok := token.Claims.(*ClaimJwt)
	if !ok {
		err = errors.New("Couldn't parse claim token")
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("Token has expired")
		return
	}
	return
}

func GetMail(sgnStr string) (email string, err error) {
	mySigningKey := []byte("AllForOne")
	token, err := jwt.ParseWithClaims(
		sgnStr,
		&ClaimJwt{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(mySigningKey), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*ClaimJwt)
	if !ok {
		err = errors.New("Couldn't parse claim token")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("Token has expired")
		return
	}
	return claims.Email, nil
}
