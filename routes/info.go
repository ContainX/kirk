package routes

import "gopkg.in/gin-gonic/gin.v1"

func Info(router *gin.Engine) {
	router.GET("/info", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "The light is green, the ship is clean",
		})
	})
}
