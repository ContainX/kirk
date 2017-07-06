package main

import (
	"github.com/ContainX/kirk/routes"
	"gopkg.in/gin-gonic/gin.v1"
)

func getRouter() *gin.Engine {

	router := gin.Default()

	router.LoadHTMLGlob("html/*")
	routes.Info(router)
	routes.Auth(router)

	return router

}
