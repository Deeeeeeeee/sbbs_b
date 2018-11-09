package common

import (
	"github.com/gin-gonic/gin"
)

func ExceptionMiddelware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if 
			}
		}
	}
}