package parse

import (
	"github.com/gin-gonic/gin"
	"github.com/muixstudio/clio/internal/aggregater/utils/parse/binding"
)

func Parse(c *gin.Context, obj any) error {

	err := c.ShouldBindBodyWithJSON(obj)
	if err != nil {
		return err
	}
	err = binding.Validate(obj)
	return err
}
