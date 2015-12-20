package main

import (
   "io/ioutil"
   "log"
   "fmt"
   "errors"

   "github.com/saromanov/notesapp/service"
   "github.com/saromanov/notesapp/messagebus"
   "github.com/saromanov/notesapp/discovery"
   "github.com/saromanov/notesapp/services"

   "github.com/hashicorp/hcl"

)

var (
	errEmptyServices = errors.New("List of services is empty")
)

// Config provides main configuration structure for running notesapp
type Config struct {
	// Name is name of project
	Name     string

	// Generatie checks for microservices
	ServiceGenChecks bool
	//ConsulConfig provides path to consul config
	Consul *discovery.Config

	// TODO: Specification for Timeout and Interval
	ConsulChecks []*discovery.Service
	// Services loads list of current services
	Services []*service.Config
	// MessageBus provides configuration for message bus(RabbitMQ)
	MessageBus *messagebus.Config
}

// LoadConfig provides loading configuration for all parts of notesapp
func LoadConfig(path string) (Config, error) {
	d, err := ioutil.ReadFile(path)
	var config Config
	if err != nil {
		return config, err
	}

	errhcl := hcl.Decode(&config, string(d))
	if errhcl != nil {
		return config, errhcl
	}
	return config, nil
}

func StartMessageBus(cfg *messagebus.Config) {

}

// StartDiscrovery provides start of service discovery with consul
func StartDiscovery(consul *discovery.Config, checks []*discovery.Service) error {
	disc, err := discovery.CreateDiscovery(consul)

	if err != nil {
		return err
	}

	for _, check := range checks {
		disc.Register(check)
	}

	return nil
}

func StartNotesapp() {

}

// getAddresses returns list of addresses of services
func getAddresses(services []*service.Config) ([]string, error) {
	if len(services) == 0 {
		return nil, errEmptyServices
	}

	result := []string{}
	for _, serviceItem := range services {
		if serviceItem != nil {
			if serviceItem.ServerAddr != "" {
				result = append(result, serviceItem.ServerAddr)
			}
		}
	}

	return result, nil
}

func main() {
	var err error
	var cfg Config
	cfg, err = LoadConfig("config.hcl")
	if err != nil {
		log.Fatal(err)
	}

	var addresses []string
	addresses, err = getAddresses(cfg.Services)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(addresses)
	fmt.Println(cfg)
}