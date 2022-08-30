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
)

// db holds database connection
var db *sql.DB

func init() {
	var err error
	pgConn := fmt.Sprintf("host=%s port=%s user=%s password=%s "+
		"dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
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
	c.IndentedJSON(http.StatusOK, allBlogs)
}
