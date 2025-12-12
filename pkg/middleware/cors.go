package middleware

import (
	"github.com/gin-gonic/gin"

	"Dash/pkg/env"
	"Dash/pkg/logger"
)

func CORS() gin.HandlerFunc {
	allowOrigin := env.GetString("CORS_ALLOW_ORIGIN")
	if allowOrigin == "" {
		logger.Warn("CORS:DEFAULT_ORIGIN_USED", map[string]interface{}{
			"warning": "CORS_ALLOW_ORIGIN not set, using default '*'",
		})
		allowOrigin = "*"
	}

	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
