package scrapy

import (
	"upper.io/db.v3/sqlite"
	"log"
	"upper.io/db.v3"
	"time"
	"strings"
)

const databaseName string = "rss_subscribe"
const rss string = "rss"
const Subscribe = "subscribe"

const collectionCreateSQL = `CREATE TABLE "` + Subscribe + `" (
	"id" INTEGER PRIMARY KEY,
	"gmt_create" DATETIME,
	"gmt_modified" DATETIME,
	"title" VARCHAR(255),
	"link" TEXT,
	"pub_date" VARCHAR(255),
	"img" TEXT,
	"description" TEXT,
	"creator" VARCHAR(255),
	"category" VARCHAR(255)
)`

func DbInit() (db.Database, error) {
	// Attempting to open database.
	sess, err := sqlite.Open(sqlite.ConnectionURL{Database: databaseName})
	if err != nil {
		return nil, err
	}

	// Collection lookup.
	col := sess.Collection(rss)
	if col.Exists() {
		return sess, nil
	}

	log.Printf("Initializing database %s...", databaseName)
	// Collection does not exists, let's create it.
	// Execute CREATE TABLE.
	if _, err = sess.Exec(collectionCreateSQL); err != nil {
		return sess, nil
	}
	return sess, nil
}

type subscribe struct {
	id int64
	gmt_create int64
	gmt_modified int64
	title string
	link string
	pub_date string
	img string
	description string
	creator string
	category string
	categorys []string
}

func IsAlreadySubcribe(storage db.Collection, url string) bool {
	result := false
	find := storage.Find().Where("link", db.Like(url))
	count, _ := find.Count()
	if count != 0 {
		result = true
	}
	return result
}

func InsertRecord(storage db.Collection, item Item) error {
	vmap := make(map[string]interface{})
	vmap["gmt_create"] =  time.Now().Unix()
	vmap["gmt_modified"] =  time.Now().Unix()
	vmap["title"] =  item.Title
	vmap["link"] =  item.Link
	vmap["pub_date"] =  item.PubDate
	vmap["img"] =  item.Img
	vmap["description"] =  item.Description
	vmap["creator"] =  item.Creator
	vmap["category"] =  strings.Join(item.Category,",")
	_, e := storage.Insert(vmap)
	return e
}
