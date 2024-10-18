package main

import (
	controllers "drm/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.MaxMultipartMemory = 100 << 20
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://193.228.90.202", "https://bookadd.ir"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.Static("/uploads", "./uploads")

	api := r.Group("/backend")
	api.POST("/upload", controllers.UploadFile)

	r.Run(":8080")
}
