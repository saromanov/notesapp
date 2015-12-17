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