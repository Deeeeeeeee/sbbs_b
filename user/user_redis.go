package user

import (
	"fmt"
	"sbbs_b/common"
)

const userJWTKey = "user:jwt:%d"

func jwt(id int64) string {
	return common.Redis().Get(fmt.Sprintf(userJWTKey, id)).String()
}
