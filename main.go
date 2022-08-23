package main

import (
	"github.com/gin-gonic/gin"

	"go-crawling/routes"
)

func main() {

	r := setupRouter()
	r.Run(":7000")
}

func setupRouter() *gin.Engine {

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	api := r.Group("/api")
	{
		api.GET("/indexing", routes.Indexing)
		api.GET("/kurs", routes.GetKurs)
		api.GET("/kurs/:symbol", routes.GetKurs)
		api.POST("/kurs", routes.PostKurs)
		api.PUT("/kurs", routes.PutKurs)
		api.DELETE("/kurs/:date", routes.DeleteKurs)
	}

	return r
}
