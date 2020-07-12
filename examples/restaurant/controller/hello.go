package controller

import "github.com/gin-gonic/gin"

type HelloController struct {
}

func (hello *HelloController) Hello(ctx *gin.Context) {
	ctx.JSON(200, map[string]interface{}{
		"message": "hello world",
	})
}

func (hell *HelloController) Router(engine *gin.Engine) {
	engine.GET("/hello", hello.Hello)
}
