package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

type postAuthor struct {
	Header  string `xml:",innerxml"`
	Cdata   string `xml:",innerxml"`
	Trailer string `xml:",innerxml"`
}

type postContent struct {
	Header  string `xml:",innerxml"`
	Cdata   string `xml:",innerxml"`
	Trailer string `xml:",innerxml"`
}

type Item struct {
	Title          string      `xml:title`
	PubDate        time.Time   `xml:pubDate`
	Creator        postAuthor  `xml:"dc:creator"`
	ContentEncoded postContent `xml:"content:encoded"`
}

type channel struct {
	Title        string `xml:"title"`
	Link         string `xml:"link"`
	Description  string `xml:"description"`
	Generator    string `xml:"generator"`
	Language     string `xml:"language"`
	WpWxrVersion string `xml:"wp:wxr_version"`
	Items        []Item `xml:"item"`
}

type DbRow struct {
	Title       string
	PostAuthor  string
	PubDate     time.Time
	PostContent string
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
	stmtOut, err := db.Prepare("SELECT post_date_gmt, post_title, (select users.user_login from wp_gpgpja_users as users where posts.post_author = users.id), post_content from wp_gpgpja_posts as posts")
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query()
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	dbRow := DbRow{}
	channel := channel{
		Title:       "Faultyvision",
		Description: "Yuhri's Blog",
		Language:    "en",
		Generator:   "https://github.com/yanfali/wpextract",
		Link:        "http://www.faultyvision.net",
	}
	var nullTime mysql.NullTime
	var content = postContent{
		Header:  "<![CDATA[",
		Trailer: "]]",
	}
	var author = postAuthor{
		Header:  "<![CDATA[",
		Trailer: "]]",
	}
	for rows.Next() {
		err := rows.Scan(&nullTime, &dbRow.Title, &dbRow.PostAuthor, &dbRow.PostContent)
		if err != nil {
			log.Fatal(err)
		}
		if nullTime.Valid {
			dbRow.PubDate = nullTime.Time
		} else {
			dbRow.PubDate = time.Now()
		}
		content.Cdata = dbRow.PostContent
		author.Cdata = dbRow.PostAuthor
		channel.Items = append(channel.Items, Item{
			Title:          dbRow.Title,
			PubDate:        dbRow.PubDate,
			ContentEncoded: content,
			Creator:        author,
		})
	}
	log.Printf("Found %d posts\n", len(channel.Items))
	enc, err := xml.MarshalIndent(channel, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	os.Stdout.Write(enc)
}
