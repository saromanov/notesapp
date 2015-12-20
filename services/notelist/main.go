package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/saromanov/notesapp/configs"
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
	if err != nil {
		logger.Error(fmt.Sprintf("%v", err))
		return
	}
	serv, err := service.CreateService(cfg)

	if err != nil {
		logger.Error(fmt.Sprintf("%v", err))
	}

	serv.HandleFunc("/api/list", func(w http.ResponseWriter, r *http.Request){
		w.Header().Add("Access-Control-Allow-Origin", "*")
		dba := serv.GetDBItem()
		//var schemas []Schema
		Error := ""
		schemas, err := dba.GetAll()
		if err != nil {
			logger.Error(fmt.Sprintf("%v", err))
			Error = fmt.Sprintf("%v", err)
		}

		ser, errser := json.Marshal(schemas)
		if errser != nil {
			logger.Error(fmt.Sprintf("%v", err))
			Error = fmt.Sprintf("%v", errser)
		}

		amqpitem := serv.GetAMQPItem()
		amqpitem.Send("notesapp", "newget")
		resp := Response{ Request: "POST", API: "insert", Info: "Added", Time: time.Now().String(),
	            Data: string(ser), Error: Error}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
        	panic(err)
    	}

	})

	logger.Info("Service list is started")
	serv.Start()
}