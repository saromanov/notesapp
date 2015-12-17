package client

import (
   "../api"
   "time"
   "errors"
)

var (
	errGetAllNotes = errors.New("Error for getting list of notes")
)

type ClientNotesapp struct {
	Addr string
}

func (cli *ClientNotesapp) CreateNote(title, noteitem string, tags[]string) error {
	var respNote api.Note
	note := api.Note(Title: title, NoteItem: noteitem, Tags: tags, 
		CreateTime: time.Now(), ModTime: time.Now())

	req, err := request(cli.Addr, "POST", node)
	if err != nil {
		return err
	}

	err := unmarshal(req, &respNote)
	if err != nil {
		return err
	}
}

func (cli *ClientNotesapp) GetAllNotes() ([]api.Note, error) {
	var respNotes []api.Note

	req, err := request(cli.Addr, "GET", nil)
	if err != nil {
		return respNotes, err
	}

	err := unmarshal(req, &respNotes)
	if err != nil {
		return err
	}

	return respNotes, err
}