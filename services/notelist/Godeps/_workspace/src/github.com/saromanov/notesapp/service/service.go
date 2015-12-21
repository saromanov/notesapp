package service

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/saromanov/notesapp/db"
	"github.com/saromanov/notesapp/publisher"
)

type (
	Handler func(w http.ResponseWriter, r *http.Request)
)

// Service provides implementation of basic Service
type Service struct {
	Title      string
	Addr       string
	Port       string
	handlers map[string]Handler
	amqp *publisher.Publisher
	dbitem     *db.DB
	running    chan bool
}

func CreateService(config *Config) (*Service, error) {
	err := CheckConfig(config) 
	if err != nil {
		return nil, err
	}

	service := new(Service)
	mongoconfig := &db.Config {
		Addr: config.MongoAddr,
		DBName: config.MongoDBName,
	}

	dbitem, err := db.CreateDB(mongoconfig)
	if err != nil {
		return nil, err
	}
	service.dbitem = dbitem
	amqp, err := publisher.NewPublisher(config.RabbitExchange, config.RabbitAddr)
	if err != nil {
		return nil, err
	}
	service.amqp = amqp
	service.handlers = map[string]Handler{}
	service.Addr = config.ServerAddr
	return service, nil

}

// HandleFunc provides append function for API
func (service *Service) HandleFunc(title string, fn Handler){
	service.handlers[title] = fn
}

// SendMessage provides sending message with RabbitMQ
func (service *Service) SendMessage(exchangename, msg string){
	service.amqp.Send(exchangename, msg)
}

// GetDBItem returns current MongoDB state
func (service *Service) GetDBItem() *db.DB {
	return service.dbitem
}

// GetDBItem returns current AMQP state
func (service *Service) GetAMQPItem() *publisher.Publisher {
	return service.amqp
}

// Start set of service is alive
func (service *Service) Start() {
	go func() {
		service.running <- true
	}()

	r := mux.NewRouter()
	for name, fn := range service.handlers {
		r.HandleFunc(name, fn)
	}

	http.ListenAndServe(service.Addr, r)
}

// Stop provides off service
func (service *Service) Stop() {
	service.amqp.Close()
	service.dbitem.Close()
}
