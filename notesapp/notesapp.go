package notesapp

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/beatrichartz/martini-sockets"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"

	"github.com/saromanov/notesapp/client"
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
}

// Add a client to a room
func (r *Room) appendClient(client *Client) {
	r.Lock()
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
	notes, err := r.getAll()
	if err != nil {
		return
	}
	encnotes, errenc := json.Marshal(notes)
	if errenc != nil {
		return
	}
	client.out <- &Note{"list", client.Name, "", string(encnotes)}
}

func (r *Room) messageOtherClients(client *Client, msg *Note) {
	r.Lock()
	msg.Title = client.Name

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

	for index, c := range r.clients {
		if c == client {
			r.clients = append(r.clients[:index], r.clients[(index+1):]...)
		} else {
			c.out <- &Note{"status", client.Name, fmt.Sprintf("%d", len(r.clients)), ""}
		}
	}
}

func (r *Room) notify(client *Client, itemname string) {
	r.Lock()
	for _, c := range r.clients {
		if c != client {
			c.out <- &Note{"checkitem", c.Name, itemname, ""}
		}
	}
	defer r.Unlock()
}

func (r *Room) insert(msg *Note) error {
	err := r.insertClient.CreateNote(msg.Title, msg.Text)
	if err != nil {
		return err
	}

	return nil
}

// getNote provides getting note by the title
func (r *Room) getNote(cli *Client, title string) (client.Schema, error) {
	return client.Schema{}, nil
}

// getAll return all notes
func (r *Room) getAll() ([]client.Schema, error) {
	var resnotes []client.Schema
	notes, err := r.getallClient.GetAllNotes()
	fmt.Println(notes)
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
		err := r.insert(msg)
		if err != nil {
			return err
		}
		r.notify(client, msg.Title)
	case "get":
		_, err := r.getNote(client, msg.Title)
		if err != nil {
			return err
		}
	}

	return nil
}

func Start() {
	m := martini.Classic()

	room := &Room{sync.Mutex{}, "test1", client.ClientNotesapp{Addr: "http://127.0.0.1:8080/api/insert"},
		client.ClientNotesapp{Addr: "http://127.0.0.1:8082/api/list"}, make([]*Client, 0)}
	// Use Renderer
	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))

	// Index
	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", "")
	})

	m.Get("/sockets/:id", sockets.JSON(Note{}), func(r render.Render, params martini.Params, receiver <-chan *Note, sender chan<- *Note, done <-chan bool, disconnect chan<- int, err <-chan error) (int, string) {
		fmt.Println(params["id"])
		client := &Client{params["id"], receiver, sender, done, err, disconnect}
		// A single select can be used to do all the messaging
		room.appendClient(client)
		for {
			select {
			case <-client.err:
				// Don't try to do this:
				// client.out <- &Message{"system", "system", "There has been an error with your connection"}
				// The socket connection is already long gone.
				// Use the error for statistics etc
			case msg := <-client.in:
				if msg.Event == "add" {
					err := room.insert(msg)
					if err != nil {
						fmt.Println(err)
					} else {
						room.notify(client, msg.Title)
					}
				} else {
					room.messageOtherClients(client, msg)
				}
				//r.messageOtherClients(client, msg)
			case <-client.done:
				fmt.Println("DONE")
				//r.removeClient(client)
				return 200, "OK"
			}
		}
	})
	m.Run()

}
