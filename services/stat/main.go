package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"log"
	"time"

	"github.com/saromanov/notesapp/configs"
	"github.com/saromanov/notesapp/db"
	"github.com/saromanov/notesapp/logging"
	"github.com/saromanov/notesapp/service"
)

type Response struct {
	Info    string
	API     string
	Time    string
	Error   string
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
	if mongo != "" && mongoport != "" {
		cfg.MongoAddr = fmt.Sprintf("%s:%s", mongo, mongoport)
	}
	rabbit := os.Getenv("RABBIT_1_PORT_5672_TCP_ADDR")
	rabbitport := os.Getenv("RABBIT_1_PORT_5672_TCP_PORT")
	if rabbit != "" && rabbitport != "" {
		cfg.RabbitAddr = fmt.Sprintf("amqp://%s:%s", rabbit, rabbitport)
	}
	time.Sleep(10 * time.Second)
	stat, err := service.CreateService(cfg)

	if err != nil {
		log.Fatal(err)
	}

	dba := stat.GetDBItem()
	dba.InsertIfNotExists("Stat", &db.DBStat{
		Title:         "Stat",
		Notes:         0,
		Starttime:     time.Now().String(),
		Getnums:       0,
		Microservices: 0,
	})

	stat.HandleFunc("/api/stat", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		dba := stat.GetDBItem()
		Error := ""
		schema, err := dba.GetStat("Stat")
		if err != nil {
			logger.Error(fmt.Sprintf("%v", err))
			Error = fmt.Sprintf("%v", err)
		}

		ser, errser := json.Marshal(schema)
		if errser != nil {
			logger.Error(fmt.Sprintf("%v", err))
			Error = fmt.Sprintf("%v", errser)
		}
		resp := Response{Request: "GET", API: "insert", Info: "Added", Time: time.Now().String(),
			Data: string(ser), Error: Error}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}

	})

	stat.HandleFunc("/api/incgets", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		dba := stat.GetDBItem()
		//var schemas []Schema
		Error := ""
		err := dba.IncGetNums("Stat")
		if err != nil {
			logger.Error(fmt.Sprintf("%v", err))
			Error = fmt.Sprintf("%v", err)
		}

		resp := Response{Request: "POST", API: "incgets", Info: "Added", Time: time.Now().String(),
			Data: "", Error: Error}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}

	})

	stat.HandleFunc("/api/incmicroservices", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		dba := stat.GetDBItem()
		//var schemas []Schema
		Error := ""
		err := dba.IncMicroservices("Stat")
		if err != nil {
			logger.Error(fmt.Sprintf("%v", err))
			Error = fmt.Sprintf("%v", err)
		}

		resp := Response{Request: "POST", API: "incmicroservices", Info: "Added", Time: time.Now().String(),
			Data: "", Error: Error}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}

	})

	stat.HandleFunc("/api/incnotes", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		dba := stat.GetDBItem()
		//var schemas []Schema
		Error := ""
		err := dba.IncNotes("Stat")
		if err != nil {
			logger.Error(fmt.Sprintf("%v", err))
			Error = fmt.Sprintf("%v", err)
		}

		resp := Response{Request: "POST", API: "incnotes", Info: "Added", Time: time.Now().String(),
			Data: "", Error: Error}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}

	})

	logger.Info("Stat service is started")
	stat.Start()
}
