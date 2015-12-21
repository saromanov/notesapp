package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/saromanov/notesapp/configs"
	"github.com/saromanov/notesapp/db"
	"github.com/saromanov/notesapp/service"
	"github.com/saromanov/notesapp/logging"
)

type Response struct {
	Info string
	API  string
	Time string
	Error string
	Request string
	Data    string
}

type InsertNoteRequest struct {
	Title string `json:"title"`
	Note  string `bson:"note"`
}


func main() {
	logger := logging.NewLogger(nil)
	args := os.Args
	if len(args) == 1 {
		logger.Error("Argument for config is not found")
		return
	}

	path := args[1]
	cfg, err := configs.LoadServiceConfig(path)
	fmt.Println(cfg)
	if err != nil {
		logger.Error(fmt.Sprintf("%v", err))
		return
	}
	mongo := os.Getenv("MONGODB_1_PORT_27017_TCP_ADDR")
	mongoport := os.Getenv("MONGODB_1_PORT_27017_TCP_PORT")
	rabbit := os.Getenv("RABBIT_1_PORT_5672_TCP_ADDR")
	rabbitport := os.Getenv("RABBIT_1_PORT_5672_TCP_PORT")
	cfg.MongoAddr = fmt.Sprintf("%s:%s", mongo, mongoport)
	cfg.RabbitAddr = fmt.Sprintf("amqp://%s:%s", rabbit, rabbitport)
	time.Sleep(10 * time.Second)
	serv, err := service.CreateService(cfg)
	if err != nil {
		logger.Error(fmt.Sprintf("%v", err))
		return
	}

	serv.HandleFunc("/api/insert", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		Error := ""
		var req InsertNoteRequest
		var err error
		err = decoder.Decode(&req)
		if err != nil {
			logger.Error(fmt.Sprintf("%s: %v", "/api/insert", err))
			Error = fmt.Sprintf("%v", err)
		}
		dba := serv.GetDBItem()
		schema := db.Schema{
			Title:      req.Title,
			NoteItem:   req.Note,
			CreateTime: time.Now(),
			ModTime:    time.Now(),
		}
		err = dba.Insert(schema)
		if err != nil {
			logger.Error(fmt.Sprintf("%s: %v", "/api/insert", err))
			Error = fmt.Sprintf("%v", err)
		}

		amqpitem := serv.GetAMQPItem()
		encschema, err := json.Marshal(schema)
		if err != nil {
			logger.Error(fmt.Sprintf("%s: %v", "/api/insert", err))
			Error = fmt.Sprintf("%v", err)
		} else {
			amqpitem.Send("notesapp", string(encschema))
		}

		logger.Info("/api/insert is complete")
		resp := Response{Request: "POST", API: "insert", Info: "Added", Time: time.Now().String(),
			Data: "", Error: Error}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}

	})
    
    logger.Info("Microservice insertnote is started")
    serv.Start()
}
