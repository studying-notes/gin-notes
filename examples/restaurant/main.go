package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"restaurant/config"
	"restaurant/controller"
)

func registerRouter(router *gin.Engine) {
	new(controller.HelloController).Router(router)
}

func main() {
	cfg, err := config.ParseConfig("config/app.json")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	app := gin.Default()
	_ = app.Run(cfg.AppHost + ":" + cfg.AppPort)
}
