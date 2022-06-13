package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
)

//Db holds database connection
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
	CheckError(err)
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
	// Routing
	router := http.NewServeMux()
	router.HandleFunc("/blog/", handleArticle)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func handleArticle(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case "GET":
		err = articleHandleGet(w, r)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
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

func articleHandleGet(w http.ResponseWriter, r *http.Request) (err error) {
	blogId, _ := strconv.Atoi(path.Base(r.URL.Path))
	blog, _ := getArticle(blogId)
	output, err := json.Marshal(blog)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {
		return err
	}
	fmt.Println(blog.Title)

	return nil
}
