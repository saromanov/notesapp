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

	"github.com/gorilla/mux"
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
	mongo := os.Getenv("MONGODB_1_PORT_27017_TCP_ADDR")
	mongoport := os.Getenv("MONGODB_1_PORT_27017_TCP_PORT")
	rabbit := os.Getenv("RABBIT_1_PORT_5672_TCP_ADDR")
	rabbitport := os.Getenv("RABBIT_1_PORT_5672_TCP_PORT")
	if(mongo != "" && mongoport != "") {
		cfg.MongoAddr = fmt.Sprintf("%s:%s", mongo, mongoport)
	}
	
	if(rabbit != "" && rabbitport != "") {
		cfg.RabbitAddr = fmt.Sprintf("amqp://%s:%s", rabbit, rabbitport)
	}

	serv, err := service.CreateService(cfg)

	if err != nil {
		logger.Error(fmt.Sprintf("%v", err))
	}

	serv.HandleFunc("/api/delete/{title}", func(w http.ResponseWriter, r *http.Request){
		vars := mux.Vars(r)
		title := vars["title"]
		dba := serv.GetDBItem()
		Error := ""
		err := dba.Remove(title)
		if err != nil {
			logger.Error(fmt.Sprintf("%v", err))
			Error = fmt.Sprintf("%v", err)
		}

		resp := Response{ Request: "POST", API: "delete", Info: "Removed", Time: time.Now().String(),
	            Data: "", Error: Error}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
        	panic(err)
    	}

	})

	logger.Info("Service deletenote is started")
	serv.Start()
}
