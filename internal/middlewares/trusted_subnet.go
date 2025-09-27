package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rfruffer/go-musthave-shortener/internal/utils"
)

// TrustedSubnetMiddleware проверяет IP адрес клиента в доверенной подсети
func TrustedSubnetMiddleware(trustedSubnet string) gin.HandlerFunc {
	return func(c *gin.Context) {
		realIP := c.GetHeader("X-Real-IP")
		if realIP == "" {
			realIP = c.ClientIP()
		}

		allowed, err := utils.IsIPInTrustedSubnet(realIP, trustedSubnet)
		if err != nil || !allowed {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
