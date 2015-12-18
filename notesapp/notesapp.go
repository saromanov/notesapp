package notesapp

import (
	"fmt"
	//"log"
	//"net/http"

	"github.com/beatrichartz/martini-sockets"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

type Note struct {
	Event string `json:"event"`
	Title string `json:"title"`
	Text  string `json: "text"`
}

type Client struct {
	Name       string
	in         <-chan *Note
	out        chan<- *Note
	done       <-chan bool
	err        <-chan error
	disconnect chan<- int
}

func Start() {
	m := martini.Classic()

	// Use Renderer
	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))

	// Index
	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", "")
	})

	// render the room
	m.Get("/list", func(r render.Render, params martini.Params) {
		r.HTML(200, "room", map[string]map[string]string{"room": map[string]string{"name": params["name"]}})
	})

	m.Get("/sockets/:id", sockets.JSON(Note{}), func(r render.Render, params martini.Params, receiver <-chan *Note, sender chan<- *Note, done <-chan bool, disconnect chan<- int, err <-chan error) (int, string) {
		client := &Client{params["id"], receiver, sender, done, err, disconnect}
		// A single select can be used to do all the messaging
		for {
			select {
			case <-client.err:
				// Don't try to do this:
				// client.out <- &Message{"system", "system", "There has been an error with your connection"}
				// The socket connection is already long gone.
				// Use the error for statistics etc
			case msg := <-client.in:
				fmt.Println("MSG: ", msg)
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
