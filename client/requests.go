package client

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"bytes"
	"net/http"
	"time"

	"github.com/saromanov/notesapp/api"
)

type Response struct {
	Info string
	API  string
	Time string
	Error string
	Request string
	Data    string
}

type Schema struct {
	//Id      bson.ObjectId `json:"_id,omitempty"`
	Title   string `json:"title"`
	NoteItem string `json: "note_item"`
	Tags    []string `json: "tags"`
	Version  int      `json: "version"`
	CreateTime  time.Time `json: "create_time"`
	ModTime     time.Time `json: "mod_time"`
}

var (
	errEmptyUrl = errors.New("Url for request is empty")
	errUndefinedMethod = errors.New("Method is undefined")
	errEmptyItem = errors.New("Item is empty")
)

type Noteslist []api.Note
// Helping function for getting requests to the server
func request(url, method string, item interface{}) (*http.Response, error) {
	if url == "" {
		return nil, errEmptyItem
	}

	if method != "GET" && method != "POST" {
		return nil, errUndefinedMethod
	}

	body, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")

	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func unmarshal(resp *http.Response, respitem interface{}) error {
	var err error
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(respBody, &respitem)
	if err != nil {
		return err
	}

	return nil
}

func unmarshalList(resp *http.Response) ([]Schema, error) {
	var err error
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var respData Response
	err = json.Unmarshal(respBody, &respData)
	if err != nil {
		return nil, err
	}

	var noteList []Schema
	err = json.Unmarshal([]byte(respData.Data), &noteList)
	if err != nil {
		return nil, err
	}

	return noteList, nil
}