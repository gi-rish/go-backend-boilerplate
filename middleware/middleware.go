package middleware

import "github.com/gin-gonic/gin"

func SampleMiddleware(c *gin.Context) {
	println("Sample Middleware")
	c.Next()
}
