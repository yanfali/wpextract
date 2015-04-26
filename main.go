package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Channel struct {
	Title        string `xml:"title"`
	Link         string `xml:"link"`
	Description  string `xml:"description"`
	Generator    string `xml:"generator"`
	Language     string `xml:"language"`
	WpWxrVersion string `xml:"wp:wxr_version"`
	Items        []Item `xml:"item"`
}

type rss struct {
	Version string   `xml:"version,attr"`
	Content string   `xml:"xmlns:content,attr"`
	Wfw     string   `xml:"xmlns:wfw,attr"`
	Dc      string   `xml:"xmlns:dc,attr"`
	Wp      string   `xml:"xmlns:wp,attr"`
	Channel *Channel `xml:"channel"`
}

type PostMetaDBRow struct {
	MetaKey   string `xml:"meta_key"`
	MetaValue string `xml:'meta_value"`
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

	channel, err := getPosts(db)
	if err != nil {
		log.Fatalf("%v", err)
	}
	rss := rss{
		Version: "2.0",
		Channel: &channel,
		Content: "http://purl.org/rss/1.0/modules/content/",
		Wfw:     "http://wellformedweb.org/CommentAPI/",
		Dc:      "http://purl.org/elements/1.1/",
		Wp:      "http://wordpress.org/export/1.0/",
	}
	log.Printf("Found %d posts\n", len(channel.Items))
	enc, err := xml.MarshalIndent(rss, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	os.Stdout.Write(enc)
}
