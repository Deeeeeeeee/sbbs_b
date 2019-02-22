package user

import (
	"encoding/base64"
	"net/http"
	"sbbs_b/common"
	"sbbs_b/entity"
	"sbbs_b/util"
	"time"

	"github.com/gin-gonic/gin"
)

const authHeader string = "Authentication"

// SetupUserAPIRouter 初始化 user api router
func SetupUserAPIRouter(r *gin.RouterGroup) {
	r.POST("/login", userLogin)
	r.POST("/register", userRegistered)
}

// userRegistered 用户注册
func userRegistered(c *gin.Context) {
	var dto entity.User
	common.BindJSONWithValidate(c, &dto, "user_register")

	// 根据邮箱查询，如果有则抛出异常，否则新增用户
	if has, _ := common.DBEngine().Table("user").Where("email = ?", dto.Email).Exist(); has {
		c.JSON(http.StatusBadRequest, gin.H{"message": "email 已经存在"})
		return
	}
	// 密码加密
	encryPwd(&dto)

	id, _ := common.DBEngine().Insert(dto)
	// 返回 id 和 jwt
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func encryPwd(user *entity.User) {
	salt, _ := util.RandomSalt(8)
	cryPwd, err := util.CryptPwd(salt, user.Password)
	if err != nil {
		panic(common.HTTP500Error("密码加密失败" + err.Error()))
	}

	user.Password = cryPwd
	user.Salt = base64.StdEncoding.EncodeToString(salt)
}

// userLogin 用户登录
func userLogin(c *gin.Context) {
	var dto entity.User
	common.BindJSONWithValidate(c, &dto, "user_login")
	// 校验密码
	var persist entity.User
	if has, _ := common.DBEngine().Where("id=?", dto.ID).Get(&persist); has != true {
		panic(common.HTTP400Error("账户不存在"))
	}
	salt, err := base64.StdEncoding.DecodeString(persist.Salt)
	if err != nil {
		panic(common.HTTP500Error("密码解密失败" + err.Error()))
	}
	if cryPwd, _ := util.CryptPwd(salt, dto.Password); cryPwd != persist.Password {
		panic(common.HTTP400Error("密码校验失败"))
	}
	// 从 redis 查询 jwt，存在直接返回，不存在创建 jwt
	jwt := jwt(dto.ID)
	if len(jwt) == 0 {
		jwt = util.GenerateJwt(persist.ID)
	}
	c.JSON(http.StatusOK, gin.H{"jwt": jwt})
}

// userLogout 用户退出登录
func userLogout(c *gin.Context) {
	// 将 jwt 放入黑名单
	token := c.GetHeader(authHeader)
	t := util.JwtExpiresAt(token) - time.Now().UnixNano()
	addBlackJwt(token, time.Duration(t)*time.Nanosecond)
}

// JwtMiddelware Jwt校验中间件
func JwtMiddelware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(authHeader)
		if ok, userID := ValidJwt(token); ok {
			c.Set("userId", userID)
			c.Next()
		} else {
			panic(common.HTTP400Error("jwt不正确或者过期"))
		}
	}
}

// ValidJwt 校验 jwt 返回是否校验通过和 userID
func ValidJwt(token string) (bool, string) {
	// 从 redis black list 对比
	if isBlack := isBlackJwt(token); isBlack {
		return false, ""
	}
	// 校验 jwt
	userID := util.ValidJwt(token)
	return true, userID
}
