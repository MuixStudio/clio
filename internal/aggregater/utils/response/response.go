package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muixstudio/clio/pkg/errorx"
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
	var hErr errorx.HttpError
	errors.As(err, &hErr)
	if hErr != nil {
		c.JSON(hErr.GetStatusCode(), gin.H{
			"code":    hErr.GetCode(),
			"message": hErr.Error(),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    -1,
		"message": err.Error(),
	})
}
