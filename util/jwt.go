package util

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// GenerateJwt 生成 jwt
func GenerateJwt(userId int64) string {
	key := []byte("lSyYRiuJyxzrrsjhPSlcmBGvTdYh")

	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().AddDate(0, 0, 1).UnixNano(),
		Issuer:    "front",
		Id:        string(userId),
		IssuedAt:  time.Now().UnixNano(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString(key)
	return ss
}
