package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"os"
	"strings"

	"github.com/beatrichartz/martini-sockets"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"

	"github.com/saromanov/notesapp/client"
	"github.com/saromanov/notesapp/logging"
)

type Note struct {
	Event string `json:"event"`
	Title string `json:"title"`
	Text  string `json: "text"`

	Items string `json: "items"`
}

type Client struct {
	Name string
	in   <-chan *Note
	out  chan<- *Note

	done       <-chan bool
	err        <-chan error
	disconnect chan<- int
}

// Room level
type Room struct {
	sync.Mutex
	name         string
	insertServiceAddr string
	getListServiceAddr string
	deleteServiceAddr string
	getNoteServiceAddr string
	clients      []*Client
	logger       *logging.Logger
}

// Add a client to a room
func (r *Room) appendClient(client *Client) {
	r.Lock()
	r.logger.Info("Append new client")
	r.clients = append(r.clients, client)
	for _, c := range r.clients {
		//if c != client {
		c.out <- &Note{"new", client.Name, fmt.Sprintf("%d", len(r.clients)), ""}
		//}
	}
	r.Unlock()
}

func (r *Room) loadNotes(client *Client) {
	r.Lock()
	defer r.Unlock()
	notes, err := r.getAll()
	if err != nil {
		return
	}

	for _, note := range notes {
		encnotes, errenc := json.Marshal(note)
		if errenc != nil {
			return
		}

		client.out <- &Note{"list", client.Name, "", string(encnotes)}
	}
}

func (r *Room) messageOtherClients(client *Client, msg *Note) {
	r.Lock()
	msg.Title = client.Name
	r.logger.Info(fmt.Sprintf("Sending new message from %s", client.Name))
	for _, c := range r.clients {
		//if c != client {
		c.out <- msg
		//}
	}
	defer r.Unlock()
}

// Remove a client from a room
func (r *Room) removeClient(client *Client) {
	r.Lock()
	defer r.Unlock()
	r.logger.Info(fmt.Sprintf("Remove client %s", client.Name))
	for index, c := range r.clients {
		if c == client {
			r.clients = append(r.clients[:index], r.clients[(index+1):]...)
		} else {
			c.out <- &Note{"status", client.Name, fmt.Sprintf("%d", len(r.clients)), ""}
		}
	}
}

func (r *Room) notify(client *Client, eventname, itemname string) {
	r.Lock()
	r.logger.Info("Notification to other clients")
	for _, c := range r.clients {
		if c != client {
			c.out <- &Note{eventname, c.Name, itemname, ""}
		}
	}
	defer r.Unlock()
}

func (r *Room) notifyUpdate(client *Client, msg *Note) {
	r.Lock()
	defer r.Unlock()
	r.logger.Info("Notification to other clients")
	for _, c := range r.clients {
		if c != client {
			c.out <- msg
		}
	}
}

func (r *Room) insert(cli *Client, msg *Note) error {
	ins := client.ClientNotesapp{Addr: fmt.Sprintf("%s/%s", r.insertServiceAddr, "api/insert")}
	err := ins.CreateNote(msg.Title, msg.Text)
	if err != nil {
		return err
	}
	return nil
}

// Update exist note
func (r *Room) updateNote(msg *Note) error {
	cli := client.ClientNotesapp{Addr: fmt.Sprintf("%s/%s", r.insertServiceAddr, "api/update")}
	return cli.UpdateNote(msg.Items, msg.Title, msg.Text)
}

func (r *Room) removeNote(msg *Note) error {
	cli := client.ClientNotesapp{Addr: fmt.Sprintf("%s/%s/%s", r.deleteServiceAddr, "api/delete/", msg.Title)}
	return cli.RemoveNote(msg.Title)
}

// getNote provides getting note by the title
func (r *Room) getNote(cli *Client, title string) (client.Schema, error) {
	c := client.ClientNotesapp{Addr: fmt.Sprintf("%s/%s/%s", r.getNoteServiceAddr, "api/get/", title)}
	return c.GetNote(title)
}

// getAll return all notes
func (r *Room) getAll() ([]client.Schema, error) {
	var resnotes []client.Schema
	cli := client.ClientNotesapp{Addr: fmt.Sprintf("%s/%s", r.getListServiceAddr, "api/list")}
	notes, err := cli.GetAllNotes()
	if err != nil {
		return resnotes, err
	}
	resnotes = notes
	return resnotes, nil
}

// processMessages provides routing messages
func (r *Room) processMessages(client *Client, msg *Note) error {
	switch msg.Event {
	case "add":
		err := r.insert(client, msg)
		if err != nil {
			return err
		}

		note, errgetting := r.getNote(client, msg.Title)
		if errgetting != nil {
			return errgetting
		}

		r.notifyUpdate(client, &Note{Event: "checkitem", Title: note.Title, Text: note.NoteItem})
		
	case "update":
		err := r.updateNote(msg)
		if err != nil {
			return err
		}
		r.notifyUpdate(client, msg)

	case "remove":
		err := r.removeNote(msg)
		if err != nil {
			return err
		}
		r.notify(client, "removeitem", msg.Title)

	default:
		r.messageOtherClients(client, msg)
	}

	return nil
}

func main() {
	m := martini.Classic()
	insertservice := os.Getenv("INSERTNOTE_PORT")
	if insertservice != ""{
		insertservice = strings.Replace(insertservice, "tcp", "http", -1)
	}

	listservice := os.Getenv("NOTELIST_PORT")
	if listservice != ""{
		listservice = strings.Replace(listservice, "tcp", "http", -1)
	}
	delservice := os.Getenv("DELETENOTE_PORT")
	if delservice != ""{
		delservice = strings.Replace(delservice, "tcp", "http", -1)
	}
	getservice := os.Getenv("GETNOTE_PORT")
	if getservice != ""{
		getservice = strings.Replace(getservice, "tcp", "http", -1)
	}
	room := &Room{sync.Mutex{}, "test1", 
	insertservice, listservice, delservice, getservice, 
	    make([]*Client, 0),
		logging.NewLogger(nil),
	}
	// Use Renderer
	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))

	// Index
	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", "")
	})

    m.Get("/sockets/:id", sockets.JSON(Note{}), func(r render.Render, params martini.Params, receiver <-chan *Note, sender chan<- *Note, done <-chan bool, disconnect chan<- int, err <-chan error) (int, string) {
		client := &Client{params["id"], receiver, sender, done, err, disconnect}
		// A single select can be used to do all the messaging
		room.loadNotes(client)
		room.appendClient(client)
		for {
			select {
			case <-client.err:
				// Don't try to do this:
				// client.out <- &Message{"system", "system", "There has been an error with your connection"}
				// The socket connection is already long gone.
				// Use the error for statistics etc
			case msg := <-client.in:
				err := room.processMessages(client, msg)
				if err != nil {
					room.logger.Error(fmt.Sprintf("%v", err))
				}
			case <-client.done:
				room.removeClient(client)
				return 200, "OK"
			}
		}
	})
	m.Run()

}
