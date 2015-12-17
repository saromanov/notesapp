package notesapp

import (
	"fmt"
	//"log"
	"net/http"

	"github.com/beatrichartz/martini-sockets"
	"github.com/go-martini/martini"
	//"github.com/googollee/go-socket.io"
	"github.com/martini-contrib/render"
)

type Note struct {
	Title    string `json:"title"`
	NoteItem string `json: "note_item"`
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

	/*server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(so socketio.Socket) {

		so.On("newnote", func(msg string) {
			fmt.Println("NEW: ", msg)
		})
	})

	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})*/

	m := martini.Classic()

	// Use Renderer
	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))

	// Index
	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", "")
	})

	m.Get("/:id", func(w http.ResponseWriter, r *http.Request, params martini.Params) {
		w.Header().Set("Content-Type", "appication/json")
		w.WriteHeader(200)
	})

	// render the room
	m.Get("/list", func(r render.Render, params martini.Params) {
		r.HTML(200, "room", map[string]map[string]string{"room": map[string]string{"name": params["name"]}})
	})

	m.Get("/socket.io/", sockets.JSON(Note{}), func(r render.Render, params martini.Params, receiver <-chan *Note, sender chan<- *Note, done <-chan bool, disconnect chan<- int, err <-chan error) (int, string) {
		client := &Client{params["id"], receiver, sender, done, err, disconnect}
		fmt.Println("AAAAAAAAAAA")
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
