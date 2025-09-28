package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muixstudio/clio/internal/aggregater/utils"
)

func WebCallAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析 access_token
		accessToken, err := c.Cookie("access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": err.Error(),
			})
			return
		}
		accessTokenClaims, err := utils.ParseAccessToken(accessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": err.Error(),
			})
			return
		}
		iss, ok := accessTokenClaims["iss"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": err.Error(),
			})
			return
		}
		c.Set("iss", iss)
		c.Next()
	}
}
