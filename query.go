package main

func queryAllBlog() ([]Article, error) {
	rows, err := db.Query("SELECT * FROM article WHERE published = true")
	if err != nil {
		return nil, err
	}
	var blogs []Article
	for rows.Next() {
		var blog Article
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

func queryBlog(id int) (Article, error) {
	blog := Article{}
	err := db.QueryRow("SELECT * FROM article WHERE id = $1", id).Scan(
		&blog.Id, &blog.CreatedDate, &blog.LastUpdate, &blog.Content, &blog.Published,
		&blog.Title, &blog.ArticleType)
	if err != nil {
		return blog, err
	}
	return blog, nil
}
