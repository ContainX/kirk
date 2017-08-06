package routes

import "github.com/gin-gonic/gin"

func Info(router *gin.Engine) {
	router.GET("/info", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "The light is green, the ship is clean",
		})
	})
}
