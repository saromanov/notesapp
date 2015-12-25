package db

import (
	"fmt"
	"log"
	"time"
	"errors"

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

	DBName string

	// Session is current session of Mongo
	Session *mgo.Session

	logger *log.Logger
}

func CreateDB(config *Config) (*DB, error) {
	if config == nil {
		return nil, errors.New("config data is empty")
	}

	db := new(DB)
	db.Addr = config.Addr
	db.DBName = config.DBName
	sess, err := mgo.Dial(db.Addr)
	if err != nil {
		return nil, fmt.Errorf("Can't connect to mongo by %s, go error %v\n", db.Addr, err)
	}
	db.Session = sess
	return db, nil
}

// Insert provides inserting of new note
func (db *DB) Insert(item interface{}) error {
	collection := db.Session.DB(db.DBName).C(Notes)
	err := collection.Insert(item)
	if err != nil {
		return err
	}

	return nil
}

// Updates provides updating item(note) in Mongo
func (db *DB) Update(oldtitle, title string, item Schema) error {
	schema := Schema{ModTime: time.Now()}
	if item.Title != "" && oldtitle != title {
		schema.Title = item.Title
	}
	if item.NoteItem != "" {
		schema.NoteItem = item.NoteItem
	}
	if len(item.Tags) > 0 {
		schema.Tags = item.Tags
	}

	fmt.Println(item)

	c := db.Session.DB(db.DBName).C(Notes)
	err := c.Update(bson.M{"title": oldtitle}, bson.M{"$set": item})
	fmt.Println(db.Get(title))
	return err
}

// Get provides getting by the title of note
func (db *DB) Get(title string) (Schema, error) {
	var schema Schema
	var err error
	c := db.Session.DB(db.DBName).C(Notes)
	err = c.Find(bson.M{"title": title}).One(&schema)
	return schema, err
}

func (db *DB) GetAll() ([]Schema, error) {
	var result []Schema
	var err error
	c := db.Session.DB(db.DBName).C(Notes)
	err = c.Find(bson.M{}).Sort("-mod_time").All(&result)
	return result, err
}

func (db *DB) Remove(title string) error {
	c := db.Session.DB(db.DBName).C(Notes)
	err := c.Remove(bson.M{"title": title})
	fmt.Println(err)
	return err

}

func (db *DB) Close() {
	db.Session.Close()
}
