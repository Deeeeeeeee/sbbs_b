package util

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func generateJwt(userId int64) {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().AddDate(0, 0, 1).UnixNano(),
		Issuer:    "front",
		Id:        string(userId),
		IssuedAt:  time.Now().UnixNano(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString([]byte("fsdfsl"))
	fmt.Println(ss)
}
