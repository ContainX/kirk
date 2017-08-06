package main

import (
	"github.com/ContainX/kirk/routes"
	"github.com/ContainX/kirk/tracer"
	"github.com/gin-gonic/gin"
)

func statsDMiddleware(c *gin.Context) {
	tracer.Get().Incr("requests.total", []string{"url:" + c.Request.URL.Path}, 1)
	c.Next()
}

func getRouter() *gin.Engine {

	router := gin.Default()

	// APM
	// t := ddTracer.NewTracerTransport(ddTracer.NewTransport(os.Getenv("HOST"), "8126"))
	// router.Use(gintrace.MiddlewareTracer("kirk", t))

	router.Use(statsDMiddleware)

	router.LoadHTMLGlob("html/*")
	router.StaticFile("/", "./static/html/index.html")
	router.Static("/assets", "./static/assets")
	routes.Info(router)
	routes.Auth(router)

	return router
}
