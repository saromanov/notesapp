package service

import (
	"net/http"
	"fmt"

	"github.com/gorilla/mux"

	"github.com/saromanov/notesapp/db"
	"github.com/saromanov/notesapp/publisher"
	"github.com/saromanov/notesapp/utils"
	"github.com/saromanov/notesapp/logging"
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
	logger     *logging.Logger
}

func CreateService(config *Config) (*Service, error) {
	var err error
	err = CheckConfig(config) 
	if err != nil {
		return nil, err
	}

	err = utils.WaitForMongo(config.MongoAddr)
	if err != nil {
		return nil, err
	}

	err = utils.WaitForRabbit(config.RabbitAddr)
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
	service.logger = logging.NewLogger(nil)
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
	r := mux.NewRouter()
	for name, fn := range service.handlers {
		r.HandleFunc(name, fn)
	}

	service.logger.Error(fmt.Sprintf("%v", http.ListenAndServe(service.Addr, r)))
}

// Stop provides off service
func (service *Service) Stop() {
	service.amqp.Close()
	service.dbitem.Close()
}
