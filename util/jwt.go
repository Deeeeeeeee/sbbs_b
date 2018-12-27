package util

import jwt "github.com/dgrijalva/jwt-go"

func generateJwt() {
	claims := &jwt.StandardClaims{
		ExpiresAt: 15000,
		Issuer:    "front",
		ID:			
	}
}