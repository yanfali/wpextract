package main

import (
	"database/sql"
	"fmt"
)

type PostMetaDbRow struct {
	MetaId    int
	PostId    int
	MetaKey   string
	MetaValue string
}

type MetaValue struct {
	Header string `xml:",innerxml"`
	Cdata  string `xml:",innerxml"`
	Footer string `xml:",innerxml"`
}

type PostMeta struct {
	MetaId    int       `xml:"wp:meta_id"`
	PostId    int       `xml:"wp:post_id"`
	MetaKey   string    `xml:"wp:meta_key"`
	MetaValue MetaValue `xml:"wp:meta_value"`
}

func getPostMetas(db *sql.DB, limit int) (postmeta []PostMeta, err error) {
	limitStmt := ""
	if limit > 0 {
		limitStmt = fmt.Sprintf(" limit %d", limit)
	}
	stmtOut, err := db.Prepare("SELECT " +
		"meta_id," +
		"post_id," +
		"meta_key," +
		"meta_value" +
		" from wp_gpgpja_postmeta" + limitStmt)
	if err != nil {
		return []PostMeta{}, err
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query()
	if err != nil {
		return []PostMeta{}, err
	}
	defer rows.Close()
	dbRow := PostMetaDbRow{}
	postMetas := []PostMeta{}
	metaValue := MetaValue{}
	for rows.Next() {
		err := rows.Scan(&dbRow.MetaId, &dbRow.PostId, &dbRow.MetaKey, &dbRow.MetaValue)
		if err != nil {
			return []PostMeta{}, err
		}
		metaValue.Cdata = dbRow.MetaValue
		postMetas = append(postMetas, PostMeta{
			MetaId:    dbRow.MetaId,
			PostId:    dbRow.PostId,
			MetaKey:   dbRow.MetaKey,
			MetaValue: metaValue,
		})

	}
	return postMetas, nil
}
