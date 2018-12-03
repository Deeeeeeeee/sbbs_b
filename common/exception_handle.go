package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ExceptionMiddelware 异常处理
func ExceptionMiddelware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if e, ok := err.(httpError); ok {
					http400ErrorHandle(c, e)
				} else {
					panic(err)
				}
			}
		}()
		c.Next()
	}
}

func http400ErrorHandle(c *gin.Context, err httpError) {
	c.AbortWithStatusJSON(err.HTTPCode(), gin.H{"message": err.Error()})
}

// HTTPError 带 http code 的 error
type httpError interface {
	Error() string
	HTTPCode() int
}

// HTTP400Error 返回 http400Error
func HTTP400Error(msg string) error {
	return &http400Error{httpCode: http.StatusBadRequest, Message: msg}
}

type http400Error struct {
	httpCode int
	Message  string
}

func (e *http400Error) Error() string {
	return e.Message
}

func (e *http400Error) HTTPCode() int {
	return e.httpCode
}
