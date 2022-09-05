package main

import (
	"github.com/dpgrahm4401/dpgraham-server/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

func routerSetup() (router *gin.Engine) {
	router = gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET"},
	}))
	return
}

func main() {
	router := routerSetup()
	router.GET("/blog", routes.GetAllBlogs)
	router.GET("/blog/:id", routes.GetBlog)
	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
