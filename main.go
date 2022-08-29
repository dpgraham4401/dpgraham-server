package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strconv"
)

// Db holds database connection
var Db *sql.DB

func init() {
	var err error
	pgConn := fmt.Sprintf("host=%s port=%s user=%s password=%s "+
		"dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWD"),
		os.Getenv("DB_NAME"))
	Db, err = sql.Open("postgres", pgConn)
	if err != nil {
		log.Fatal(err)
	}
}

// Article captures metadata about a blog post or tutorial etc.
type Article struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	LastUpdate  string `json:"updateDate"`
	CreateDate  string `json:"createDate"`
	Published   bool   `json:"publish"`
	ArticleType string `json:"type"`
	Content     string `json:"content"`
}

func main() {
	router := gin.Default()
	router.GET("/blog/:id", getBlogs)
	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

func getBlogs(c *gin.Context) {
	id := c.Param("id")
	idInt, _ := strconv.Atoi(id)
	myArticle, _ := getArticle(idInt)
	c.IndentedJSON(http.StatusOK, myArticle)
}

func getArticle(id int) (Article, error) {
	blog := Article{}
	err := Db.QueryRow("SELECT * FROM article WHERE id = $1", id).Scan(
		&blog.Id, &blog.Title, &blog.LastUpdate, &blog.Published,
		&blog.ArticleType, &blog.CreateDate, &blog.Content)
	if err != nil {
		return blog, err
	}
	return blog, nil
}
