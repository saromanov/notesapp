package client

import (
   "time"
   "errors"
   "net/http"
   "../api"
)


var (
	errGetAllNotes = errors.New("Error for getting list of notes")
)

type ClientNotesapp struct {
	Addr string
}

func (cli *ClientNotesapp) CreateNote(title, noteitem string) error {
	var err error
	var respNote api.Note
	note := api.Note{Title: title, NoteItem: noteitem, 
		CreateTime: time.Now(), ModTime: time.Now()}

	req, err := request(cli.Addr, "POST", note)
	if err != nil {
		return err
	}

	err = unmarshal(req, &respNote)
	if err != nil {
		return err
	}

	return nil
}

func (cli *ClientNotesapp) GetAllNotes() ([]Schema, error) {
	var err error
	var req *http.Response
	req, err = request(cli.Addr, "GET", nil)
	if err != nil {
		return nil, err
	}

	respresult, errs := unmarshalList(req)
	if errs != nil {
		return nil, err
	}

	return respresult, err
}