package db

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	NotesApp = "notesapp"
	Notes    = "notes"
)

type DB struct {
	// Addr is address to mongodb
	Addr string

	// Session is current session of Mongo
	Session *mgo.Session

	logger *log.Logger
}

func CreateDB(addr string) (*DB, error) {
	db := new(DB)
	db.Addr = addr
	sess, err := mgo.Dial(addr)
	if err != nil {
		return nil, fmt.Errorf("Can't connect to mongo, go error %v\n", err)
	}
	db.Session = sess
	return db, nil
}

// Insert provides inserting of new note
func (db *DB) Insert(item interface{}) error {
	collection := db.Session.DB(NotesApp).C(Notes)
	err := collection.Insert(item)
	if err != nil {
		return err
	}

	return nil
}

// Updates provides updating item(note) in Mongo
func (db *DB) Update(id string, item Schema) error {
	schema := Schema{ModTime: time.Now()}
	if item.Title != "" {
		schema.Title = item.Title
	}
	if item.NoteItem != "" {
		schema.NoteItem = item.Title
	}
	if len(item.Tags) > 0 {
		schema.Tags = item.Tags
	}

	c := db.Session.DB(NotesApp).C(Notes)
	idhex := bson.ObjectIdHex(id)
	return c.Update(bson.M{"_id": idhex}, 
		bson.M{"$set": bson.M{"title": item.Title}, "$inc": bson.M{"version": 1}})
}

// Get provides getting by the title of note
func (db *DB) Get(title string) (Schema, error) {
	var schema Schema
	var err error
	c := db.Session.DB(NotesApp).C(Notes)
	err = c.Find(bson.M{"title": title}).One(&schema)
	return schema, err
}

func (db *DB) GetAll() ([]Schema, error) {
	var result []Schema
	var err error
	c := db.Session.DB(NotesApp).C(Notes)
	err = c.Find(bson.M{}).Sort("-mod_time").All(&result)
	return result, err
}

func (db *DB) Remove(title string) error {
	c := db.Session.DB(NotesApp).C(Notes)
	err := c.Remove(bson.M{"title": title})
	return err

}

func (db *DB) Close() {
	db.Session.Close()
}
