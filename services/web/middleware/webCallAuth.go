package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muixstudio/clio/services/web/utils/jwt"
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
			c.Abort()
			return
		}

		// Parse access_token
		accessTokenClaims, err := jwt.ParseAccessToken(accessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    10001,
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		userId, ok := accessTokenClaims["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    10002,
				"message": "internal server error",
			})
			c.Abort()
			return
		}
		c.Set("user_id", uint64(userId))

		iss, ok := accessTokenClaims["iss"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    10003,
				"message": "internal server error",
			})
			c.Abort()
			return
		}
		c.Set("iss", iss)
		c.Next()
	}
}
