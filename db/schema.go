package db

import (
    "time"

	"gopkg.in/mgo.v2/bson"
)

type Schema struct {
	Id      bson.ObjectId `bson:"_id,omitempty"`
	Title   string `bson:"title"`
	NoteItem string `bson: "note_item"`
	Tags    []string `bson: "tags"`
	Version  int      `bson: "version"`
	CreateTime  time.Time `bson: "create_time"`
	ModTime     time.Time `bson: "mod_time"`
}

// DBStat provides information about statistics in service
type DBStat struct {
	Id      bson.ObjectId `bson:"_id,omitempty"`
	Title  string `bson:"title"`
	Notes int `bson:"notes"`
	Starttime string `bson:"start_time"`
	Getnums  int `bson:"get_nums"`
	Microservices int `bson:"microservices"`
}