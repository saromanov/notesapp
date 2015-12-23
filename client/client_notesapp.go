package client

import (
   "errors"
   "net/http"
   "github.com/saromanov/notesapp/api"
)


var (
	errGetAllNotes = errors.New("Error for getting list of notes")
)

type ClientNotesapp struct {
	Addr string
}

type InsertNoteRequest struct {
	Title string `json:"title"`
	Note  string `json:"note"`
}

type RemoveNoteRequest struct {
	Title string `json:"title"`
}

type UpdateNoteRequest struct {
	OldTitle string `json:"old_title"`
	Title string `json:"title"`
	Note  string `json:"note"`
}

// CreateNote provides sending request for update note
func (cli *ClientNotesapp) CreateNote(title, noteitem string) error {
	var err error
	var respNote api.Note
	note := InsertNoteRequest{Title: title, Note: noteitem}

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

// GetAllNotes provides sending request for getting all notes
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

// Remove Note provides sending request for remove note
func (cli *ClientNotesapp) RemoveNote(title string) error {
	var err error
	var respNote api.Note
	//note := RemoveNoteRequest{Title: title}

	req, err := request(cli.Addr, "GET", nil)
	if err != nil {
		return err
	}

	err = unmarshal(req, &respNote)
	if err != nil {
		return err
	}

	return nil
}

// UpdateNote provides sending request for update note
func (cli *ClientNotesapp) UpdateNote(oldtitle, title, noteitem string) error {
	var err error
	var respNote api.Note
	note := UpdateNoteRequest{OldTitle:oldtitle, Title: title, Note: noteitem}

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