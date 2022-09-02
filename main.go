package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// db holds database connection
var db *sql.DB

func init() {
	var err error
	pgConn := fmt.Sprintf("host=%s port=%s user=%s password=%s "+
		"dbname=%s sslmode=disable",
		os.Getenv("DPG_DB_HOST"),
		os.Getenv("DPG_DB_PORT"),
		os.Getenv("DPG_DB_USER"),
		os.Getenv("DPG_DB_PASSWORD"),
		os.Getenv("DPG_DB_NAME"))
	db, err = sql.Open("postgres", pgConn)
	if err != nil {
		log.Fatal(err)
	}
}

// Article captures metadata about a blog post or tutorial etc.
type Article struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	LastUpdate  string `json:"updateDate"`
	CreatedDate string `json:"createDate"`
	Published   bool   `json:"publish"`
	ArticleType string `json:"type"`
	Content     string `json:"content"`
}

func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET"},
	}))
	router.GET("/blog", getAllBlogs)
	router.GET("/blog/:id", getBlog)
	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

func getBlog(c *gin.Context) {
	id := c.Param("id")
	idInt, _ := strconv.Atoi(id)
	blogById, _ := queryBlog(idInt)
	c.IndentedJSON(http.StatusOK, blogById)
}

func getAllBlogs(c *gin.Context) {
	allBlogs, _ := queryAllBlog()
	for i := 0; i < len(allBlogs); i++ {
		// Todo: expand on markdown to text
		//  either bring in a third party package, or write something yourself
		allBlogs[i].Content = strings.ReplaceAll(allBlogs[i].Content, "#", "")
		allBlogs[i].Content = strings.ReplaceAll(allBlogs[i].Content, "~", "")
		if len(allBlogs[i].Content) >= 200 {
			allBlogs[i].Content = allBlogs[i].Content[0:200]
		}
	}
	c.IndentedJSON(http.StatusOK, allBlogs)
}
