package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	kratosErrors "github.com/go-kratos/kratos/v2/errors"
)

func SuccessOK(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "ok",
	})
}

func SuccessWithData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "ok",
		"data":    data,
	})
}

func FailH(c *gin.Context, err error) {
	hErr := kratosErrors.FromError(err)
	m := make(map[string]string)
	c.JSON(int(hErr.GetCode()), gin.H{
		"code":     hErr.GetCode(),
		"reason":   hErr.GetReason(),
		"message":  hErr.GetMessage(),
		"metadata": &m,
	})
}
