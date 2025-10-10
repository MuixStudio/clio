package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SuccessOK(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
	})
}

func SuccessWithData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data":    data,
	})
}

func FailH(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    -1,
		"message": err.Error(),
	})
}
