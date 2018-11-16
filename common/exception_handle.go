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
				if e, ok := err.(ParseError); ok {
					parseErrorHandle(c, e)
				} else {
					panic(err)
				}
			}
		}()
		c.Next()
	}
}

func parseErrorHandle(c *gin.Context, err ParseError) {
	c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误: " + err.Error()})
}

// ParseError 解析异常
type ParseError struct {
	Message string
}

func (e *ParseError) Error() string {
	return e.Message
}
