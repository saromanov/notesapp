package api

import (
   "time"
)


type Note struct {
	Id      int32 `json:"id"`
	Title   string `json:"title"`
	NoteItem string `json: "note_item"`
	Version  int     `json: "version"`
	CreateTime  time.Time `json: "create_time"`
	ModTime     time.Time `json: "mod_time"`
}
