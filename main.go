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
	myArticle, _ := queryBlog(idInt)
	c.IndentedJSON(http.StatusOK, myArticle)
}

func getAllBlogs(c *gin.Context) {
	myArticle, _ := queryAllBlog()
	c.IndentedJSON(http.StatusOK, myArticle)
}

func queryAllBlog() ([]Article, error) {
	rows, err := db.Query("SELECT * FROM article WHERE published = true")
	if err != nil {
		return nil, err
	}
	var blogs []Article
	for rows.Next() {
		var blog Article
		if err := rows.Scan(&blog.Id, &blog.Title, &blog.LastUpdate, &blog.Published,
			&blog.ArticleType, &blog.Content, &blog.CreatedDate); err != nil {
			return blogs, err
		}
		blogs = append(blogs, blog)
	}
	if err != nil {
		return nil, err
	}
	return blogs, nil
}

func queryBlog(id int) (Article, error) {
	blog := Article{}
	err := db.QueryRow("SELECT * FROM article WHERE id = $1", id).Scan(
		&blog.Id, &blog.Title, &blog.LastUpdate, &blog.Published,
		&blog.ArticleType, &blog.Content, &blog.CreatedDate)
	if err != nil {
		return blog, err
	}
	return blog, nil
}
