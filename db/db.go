package db

import (
	"database/sql"
	"fmt"
	"github.com/dpgrahm4401/dpgraham-server/models"
	_ "github.com/lib/pq"
	"log"
	"os"
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

func QueryAllBlog() ([]models.Article, error) {
	allBlogQuery, err := db.Prepare(
		` SELECT
                     id,
                     to_char(created_date, 'mm/dd/yyyy'),
                     to_char(updated_date, 'mm/dd/yyyy'),
                     article_content,
                     published,
                     title,
                     type
                FROM
                    article
                WHERE
                    published = true`)
	rows, err := allBlogQuery.Query()
	if err != nil {
		return nil, err
	}
	var blogs []models.Article
	for rows.Next() {
		var blog models.Article
		if err := rows.Scan(&blog.Id, &blog.CreatedDate, &blog.LastUpdate, &blog.Content, &blog.Published,
			&blog.Title, &blog.ArticleType); err != nil {
			return blogs, err
		}
		blogs = append(blogs, blog)
	}
	if err != nil {
		return nil, err
	}
	return blogs, nil
}

func QueryBlog(id int) (models.Article, error) {
	blog := models.Article{}
	err := db.QueryRow("SELECT * FROM article WHERE id = $1", id).Scan(
		&blog.Id, &blog.CreatedDate, &blog.LastUpdate, &blog.Content, &blog.Published,
		&blog.Title, &blog.ArticleType)
	if err != nil {
		return blog, err
	}
	return blog, nil
}
