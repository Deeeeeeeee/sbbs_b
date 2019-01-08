package util

import (
	"fmt"
	"sbbs_b/common"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var key = []byte("lSyYRiuJyxzrrsjhPSlcmBGvTdYh")

// GenerateJwt 生成 jwt
func GenerateJwt(userId int64) string {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().AddDate(0, 0, 1).UnixNano(),
		Issuer:    "front",
		Id:        string(userId),
		IssuedAt:  time.Now().UnixNano(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	result, _ := token.SignedString(key)
	return result
}

// ValidJwt jwt 校验
func ValidJwt(tokenStr string) string {
	token, _ := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})

	if claims, ok := token.Claims.(jwt.StandardClaims); ok && token.Valid {
		return claims.Id
	}

	panic(common.HTTP400Error("登录过期"))
}
