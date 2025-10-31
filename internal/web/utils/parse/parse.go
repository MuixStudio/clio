package parse

import (
	"github.com/gin-gonic/gin"
	kratosErrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/muixstudio/clio/internal/web/utils/parse/binding"
)

func Parse(c *gin.Context, obj any) error {

	err := c.ShouldBindBodyWithJSON(obj)
	if err != nil {
		return kratosErrors.InternalServer("INTERNAL_SERVER_ERROR", "internal server error").WithCause(err)
	}
	err = binding.Validate(obj)
	if err != nil {
		return kratosErrors.BadRequest("BAD_REQUEST", err.Error()).WithCause(err)
	}
	return nil
}
