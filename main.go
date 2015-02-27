package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

type Post struct {
	PostedDate mysql.NullTime `json:post_date`
	Title      string         `json:post_title`
	Excerpt    string         `json:post_excerpt`
	Content    string         `json:post_content`
}

func main() {
	db, err := sql.Open("mysql", "root:@/wp")
	if err != nil {
		log.Fatalf("%v", err)
	}

	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatalf("%v", err)
	}

	stmtOut, err := db.Prepare("SELECT post_date_gmt, post_title, post_excerpt, post_content from wp_gpgpja_posts")
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query()
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	post := Post{}
	posts := []Post{}
	for rows.Next() {
		err := rows.Scan(&post.PostedDate, &post.Title, &post.Excerpt, &post.Content)
		if err != nil {
			log.Fatal(err)
		}
		posts = append(posts, post)
	}
	enc := json.NewEncoder(os.Stdout)
	log.Printf("Found %d posts\n", len(posts))
	enc.Encode(posts)
}
