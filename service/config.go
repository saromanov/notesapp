package service

import (
	"errors"
	"io/ioutil"

	"github.com/hashicorp/hcl"
)

var (
	errNilPointer          = errors.New("Config struct is empty")
	errEmptyName           = errors.New("Name must be non-empty")
	errEmptyRabbitAddr     = errors.New("RabbitAddr must be non-empty")
	errEmptyRabbitExchange = errors.New("RabbitExchange must be non-empty")
	errEmptyServer         = errors.New("Server must be non-empty")
	errEmptyMongoAddr      = errors.New("MongoAddr must be non-empty")
)

// Config provides configuration for service
type Config struct {
	Name           string
	RabbitAddr     string
	RabbitExchange string
	ServerAddr     string
	MongoAddr      string
}

// CheckConfig returns error if contains some blank fields
func CheckConfig(config *Config) error {
	if config == nil {
		return errNilPointer
	}
	if config.Name == "" {
		return errEmptyName
	}

	if config.RabbitAddr == "" {
		return errEmptyRabbitAddr
	}

	if config.RabbitExchange == "" {
		return errEmptyRabbitExchange
	}

	if config.MongoAddr == "" {
		return errEmptyMongoAddr
	}

	if config.ServerAddr == "" {
		return errEmptyServer
	}

	return nil
}

// LoadConfig provides loading configuration from config .hcl file
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
