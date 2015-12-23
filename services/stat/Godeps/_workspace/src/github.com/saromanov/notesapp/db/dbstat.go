package db

import (
"gopkg.in/mgo.v2/bson"
)

// Insert of NotExist provides inserting data to DB if title is not exist
func (db *DB) InsertIfNotExists(title string, item interface{}) error {
	var schema DBStat
	var err error
	c := db.Session.DB(db.DBName).C(Notes)
	err = c.Find(bson.M{"title": title}).One(&schema)
	if err != nil {
		return db.Insert(item)
	}

	return nil
}

// Get provides getting information about statictics
func (db *DB) GetStat(title string) (DBStat, error) {
	var stat DBStat
	var err error
	c := db.Session.DB(db.DBName).C(Notes)
	err = c.Find(bson.M{"title": title}).One(&stat)
	return stat, err
}

func (db *DB) IncNotes(id string) error {
	return db.inc(id, "notes")
}

func (db *DB) IncMicroservices(id string) error {
	return db.inc(id, "microservices")
}

func (db *DB) IncGetNums(id string) error {
	return db.inc(id, "get_nums")
}

func (db *DB) inc(id, field string) error {
	c := db.Session.DB(db.DBName).C(Notes)
	//idhex := bson.ObjectIdHex(id)
	return c.Update(bson.M{"title": id}, bson.M{"$inc": bson.M{field: 1}})
}
