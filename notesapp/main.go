package main

import (
	"encoding/json"
	"fmt"
	"sync"

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
	insertClient client.ClientNotesapp
	getallClient client.ClientNotesapp
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
	err := r.insertClient.CreateNote(msg.Title, msg.Text)
	if err != nil {
		return err
	}

	r.notify(cli, "checkitem", msg.Title)
	return nil
}

// Update exist note
func (r *Room) updateNote(msg *Note) error {
	cli := client.ClientNotesapp{Addr: "http://127.0.0.1:8080/api/update"}
	return cli.UpdateNote(msg.Items, msg.Title, msg.Text)
}

func (r *Room) removeNote(msg *Note) error {
	cli := client.ClientNotesapp{Addr: fmt.Sprintf("%s%s", "http://127.0.0.1:8084/api/delete/", msg.Title)}
	return cli.RemoveNote(msg.Title)
}

// getNote provides getting note by the title
func (r *Room) getNote(cli *Client, title string) (client.Schema, error) {
	return client.Schema{}, nil
}

// getAll return all notes
func (r *Room) getAll() ([]client.Schema, error) {
	var resnotes []client.Schema
	notes, err := r.getallClient.GetAllNotes()
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
		r.notify(client, "checkitem", msg.Title)

	case "get":
		_, err := r.getNote(client, msg.Title)
		if err != nil {
			return err
		}
		r.notify(client, "checkitem", msg.Title)

	case "update":
		err := r.updateNote(msg)
		if err != nil {
			return err
		}
		fmt.Println(msg)
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

	room := &Room{sync.Mutex{}, "test1", client.ClientNotesapp{Addr: "http://127.0.0.1:8080/api/insert"},
		client.ClientNotesapp{Addr: "http://127.0.0.1:8082/api/list"}, make([]*Client, 0),
		logging.NewLogger(nil),}
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
