package discovery

import (
  "time"
)


// Config for Consul initialization(CreateDiscovery)
type Config struct {
	Address string
	Token  string
}


type Service struct {
	ID   string
	Name string
	Tags []string
	Port  int
	Address string

	Checks []*Check
}

type Check struct {
	ID     string
	Name   string
	Address string
	Interval time.Duration
	Timeout  time.Duration
}
