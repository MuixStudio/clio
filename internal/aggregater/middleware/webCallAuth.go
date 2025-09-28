package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muixstudio/clio/internal/aggregater/utils"
)

func WebCallAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get access_token from Cookie
		accessToken, err := c.Cookie("access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    10001,
				"message": "Unauthorized",
			})
			return
		}

		// Parse access_token
		accessTokenClaims, err := utils.ParseAccessToken(accessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    10001,
				"message": "Unauthorized",
			})
			return
		}

		userId, ok := accessTokenClaims["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    10002,
				"message": "internal server error",
			})
			return
		}
		c.Set("user_id", uint32(userId))

		iss, ok := accessTokenClaims["iss"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    10003,
				"message": "internal server error",
			})
			return
		}
		c.Set("iss", iss)
		c.Next()
	}
}
