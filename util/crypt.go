package util

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/scrypt"
)

// RandomSalt 返回 8 byte 随机盐
func RandomSalt(len int) (salt []byte, err error) {
	b := make([]byte, len)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}

// CryptPwd 返回 32 位字符 scrypt 加密结果
func CryptPwd(salt []byte, pwd string) (encode string, err error) {
	dk, err := scrypt.Key([]byte(pwd), salt, 1<<15, 8, 1, 32)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(dk), nil
}
