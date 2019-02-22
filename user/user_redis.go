package user

import (
	"fmt"
	"sbbs_b/common"
	"time"
)

const (
	userJWTKey      = "user:jwt:%d"
	userJWTBlackKey = "user:jwt:black:%s"
)

func jwt(id int64) string {
	return common.Redis().Get(fmt.Sprintf(userJWTKey, id)).String()
}

func isBlackJwt(token string) bool {
	if value := common.Redis().Get(fmt.Sprintf(userJWTBlackKey, token)).String(); len(value) == 0 {
		return false
	}
	return true
}

func addBlackJwt(token string, t time.Duration) {
	common.Redis().Set(fmt.Sprintf(userJWTBlackKey, token), "1", t)
}
