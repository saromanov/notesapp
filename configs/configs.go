package configs

import (
   "io/ioutil"

   "github.com/saromanov/notesapp/service"
   "github.com/saromanov/notesapp/messagebus"

   "github.com/hashicorp/hcl"

)

// LoadServiceConfig provides loading configurations for service
func LoadServiceConfig(path string) (*service.Config, error) {
	d, err := ioutil.ReadFile(path)
	var config *service.Config
	if err != nil {
		return config, err
	}

	errhcl := hcl.Decode(&config, string(d))
	if errhcl != nil {
		return config, errhcl
	}
	return config, nil
}

// LoadMessageBusConfig loas configuration for Rabbitmq
func LoadMessageBusConfig(path string)(messagebus.Config, error) {
	d, err := ioutil.ReadFile(path)
	var config messagebus.Config
	if err != nil {
		return config, err
	}

	errhcl := hcl.Decode(&config, string(d))
	if errhcl != nil {
		return config, errhcl
	}
	return config, nil
}