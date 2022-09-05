package routes

import (
	"github.com/dpgrahm4401/dpgraham-server/db"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func GetBlog(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	blogById, err := db.QueryBlog(idInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	} else {
		c.IndentedJSON(http.StatusOK, blogById)
	}
}

func GetAllBlogs(c *gin.Context) {
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
