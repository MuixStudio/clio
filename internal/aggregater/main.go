package main

import (
	"github.com/gin-gonic/gin"
	"github.com/muixstudio/clio/internal/aggregater/handler"
)

func main() {

	r := gin.Default()

	handler.Register(&r.RouterGroup)

	if err := r.Run(":5020"); err != nil {
		panic(err)
	}
}
