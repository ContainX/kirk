package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"github.com/jeremyroberts0/kirk/routes"
)

func getRouter() *gin.Engine{

	router := gin.Default()

	router.LoadHTMLGlob("html/*")
	routes.Info(router)
	routes.Auth(router)


	return router

}
