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
					httpErrorHandle(c, e)
				} else {
					panic(err)
				}
			}
		}()
		c.Next()
	}
}

func httpErrorHandle(c *gin.Context, err httpError) {
	c.AbortWithStatusJSON(err.HTTPCode(), gin.H{"message": err.Error()})
}

// HTTPError 返回 httpError
func HTTPError(httpCode int, msg string) error {
	return &httpError{httpCode: httpCode, Message: msg}
}

// HTTP400Error 返回 httpError 其中 httpCode 为 400
func HTTP400Error(msg string) error {
	return &httpError{httpCode: http.StatusBadRequest, Message: msg}
}

// HTTP500Error 返回 httpError 其中 httpCode 为 500
func HTTP500Error(msg string) error {
	return &httpError{httpCode: http.StatusInternalServerError, Message: msg}
}

type httpError struct {
	httpCode int
	Message  string
}

func (e *httpError) Error() string {
	return e.Message
}

func (e *httpError) HTTPCode() int {
	return e.httpCode
}
