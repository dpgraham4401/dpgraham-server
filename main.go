package main

import (
	"github.com/dpgrahm4401/dpgraham-server/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func getBlog(c *gin.Context) {
	id := c.Param("id")
	idInt, _ := strconv.Atoi(id)
	blogById, err := db.QueryBlog(idInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "")
	} else {
		c.IndentedJSON(http.StatusOK, blogById)
	}
}

func getAllBlogs(c *gin.Context) {
	allBlogs, err := db.QueryAllBlog()
	if err != nil {
		c.JSONP(http.StatusInternalServerError, "")
	}
	for i := 0; i < len(allBlogs); i++ {
		allBlogs[i].Content = strings.ReplaceAll(allBlogs[i].Content, "#", "")
		allBlogs[i].Content = strings.ReplaceAll(allBlogs[i].Content, "~", "")
		if len(allBlogs[i].Content) >= 200 {
			allBlogs[i].Content = allBlogs[i].Content[0:200]
		}
	}
	c.IndentedJSON(http.StatusOK, allBlogs)
}

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
	router.GET("/blog", getAllBlogs)
	router.GET("/blog/:id", getBlog)
	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
