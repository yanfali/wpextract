package main

import (
	"database/sql"
	"fmt"
	"log"
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

type postExcerpt struct {
	Header  string `xml:",innerxml"`
	Cdata   string `xml:",innerxml"`
	Trailer string `xml:",innerxml"`
}

type Item struct {
	Title          string      `xml:"title"`
	PubDate        time.Time   `xml:"pubDate"`
	Creator        postAuthor  `xml:"dc:creator"`
	Guid           Guid        `xml:"guid"`
	ContentEncoded postContent `xml:"content:encoded"`
	ExcerptEncoded postExcerpt `xml:"excerpt:encoded"`
	PostId         int         `xml:"wp:post_id"`
	PostDate       time.Time   `xml:"wp:post_date"`
	PostDateGMT    time.Time   `xml:"wp:post_date_gmt"`
	CommentStatus  string      `xml:"wp:comment_status"`
	PingStatus     string      `xml:"wp:ping_status"`
	PostName       string      `xml:"wp:post_name"`
	PostStatus     string      `xml:"wp:status"`
	PostParentId   int         `xml:"wp:post_parent"`
	MenuOrder      int         `xml:"wp:menu_order"`
	PostType       string      `xml:"wp:post_type"`
	PostPassword   string      `xml:"wp:post_password"`
}

type Guid struct {
	Permalink bool   `xml:"isPermaLink,attr"`
	Content   string `xml:",innerxml"`
}

type PostDbRow struct {
	Title         string
	PostAuthor    string
	PubDate       time.Time
	PostDate      time.Time
	PostDateGMT   time.Time
	PostContent   string
	PostId        int
	PostParentId  int
	PostName      string
	MenuOrder     int
	PostStatus    string
	CommentStatus string
	PingStatus    string
	PostType      string
	PostPassword  string
	Guid          Guid
	PostExcerpt   string
}

func getPosts(db *sql.DB, limit int) (channel *Channel, err error) {
	limitStmt := ""
	if limit > 0 {
		limitStmt = fmt.Sprintf(" limit %d", limit)
	}
	stmtOut, err := db.Prepare("SELECT " +
		"post_date_gmt," +
		"post_date," +
		"post_title," +
		"(select users.user_login from wp_gpgpja_users as users where posts.post_author = users.id)," +
		"post_content," +
		"post_parent," +
		"id," +
		"post_name," +
		"menu_order," +
		"post_status," +
		"comment_status," +
		"ping_status," +
		"post_type," +
		"post_password," +
		"guid, " +
		"post_excerpt" +
		" from wp_gpgpja_posts as posts " + limitStmt)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query()
	if err != nil {
		return &Channel{}, err
	}

	defer rows.Close()
	dbRow := PostDbRow{}
	channel = &Channel{
		Title:       "Faultyvision",
		Description: "Yuhri's Blog",
		Language:    "en",
		Generator:   "https://github.com/yanfali/wpextract",
		Link:        "http://www.faultyvision.net",
	}
	var postTimeGMT mysql.NullTime
	var postTime mysql.NullTime
	var content = postContent{
		Header:  "<![CDATA[",
		Trailer: "]]",
	}
	var excerpt = postExcerpt{
		Header:  "<![CDATA[",
		Trailer: "]]",
	}
	var author = postAuthor{
		Header:  "<![CDATA[",
		Trailer: "]]",
	}
	for rows.Next() {
		err := rows.Scan(
			&postTimeGMT,
			&postTime,
			&dbRow.Title,
			&dbRow.PostAuthor,
			&dbRow.PostContent,
			&dbRow.PostId,
			&dbRow.PostParentId,
			&dbRow.PostName,
			&dbRow.MenuOrder,
			&dbRow.CommentStatus,
			&dbRow.PostStatus,
			&dbRow.PingStatus,
			&dbRow.PostType,
			&dbRow.PostPassword,
			&dbRow.Guid.Content,
			&dbRow.PostExcerpt,
		)
		if err != nil {
			return &Channel{}, err
		}
		if postTimeGMT.Valid {
			dbRow.PostDateGMT = postTimeGMT.Time
		} else {
			dbRow.PostDateGMT = time.Now()
		}
		if postTime.Valid {
			dbRow.PubDate = postTime.Time
			dbRow.PostDate = postTime.Time
		} else {
			dbRow.PubDate = time.Now()
			dbRow.PostDate = time.Now()
		}
		content.Cdata = dbRow.PostContent
		excerpt.Cdata = dbRow.PostExcerpt
		author.Cdata = dbRow.PostAuthor
		item := Item{
			Title:          dbRow.Title,
			PubDate:        dbRow.PubDate,
			ContentEncoded: content,
			ExcerptEncoded: excerpt,
			Creator:        author,
			PostId:         dbRow.PostId,
			PostParentId:   dbRow.PostParentId,
			PostName:       dbRow.PostName,
			MenuOrder:      dbRow.MenuOrder,
			PostStatus:     dbRow.PostStatus,
			CommentStatus:  dbRow.CommentStatus,
			PingStatus:     dbRow.PingStatus,
			PostType:       dbRow.PostType,
			PostPassword:   dbRow.PostPassword,
		}
		item.Guid = Guid{
			Permalink: false,
			Content:   dbRow.Guid.Content,
		}
		channel.Items = append(channel.Items, item)
	}
	return channel, nil
}
