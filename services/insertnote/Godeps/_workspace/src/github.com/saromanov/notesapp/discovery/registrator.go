package discovery

import (
	consul "github.com/hashicorp/consul/api"
)

type notesappAPI interface {
	ServiceRegister(*consul.AgentServiceRegistration)
	CheckRegister(*consul.AgentCheckRegistration)
	Services() (map[string]*consul.AgentService, error)
	Checks() (map[string]*consul.AgentCheck, error)
}

// NotesAppService provides struct for service registration
type NotesAppService struct {
	client *consul.Client
}

// CreateDiscovery provides initialization of Consul
func CreateDiscovery(cfg *Config) (*NotesAppService, error) {
	var cclient *consul.Client
	var err error
	consulConfig := consul.DefaultConfig()
	consulConfig.Address = cfg.Address
	if cclient, err = consul.NewClient(consulConfig); err != nil {
		return nil, err
	}

	return &NotesAppService{client: cclient}, err

}

// Register provides registration of service
func (nas *NotesAppService) Register(service *Service) error {
	reg := &consul.AgentServiceRegistration{
		ID:      service.ID,
		Name:    service.Name,
		Tags:    service.Tags,
		Port:    service.Port,
		Address: service.Address,
	}

	if err := nas.client.Agent().ServiceRegister(reg); err != nil {
		return err
	}

	for _, check := range service.Checks {
		err := nas.RegisterCheck(check)
		if err != nil {
			return err
		}
	}

	return nil
}

// RegisterCheck provides registration of check for service
func (mas *NotesAppService) RegisterCheck(check *Check) error {
	reg := &consul.AgentCheckRegistration{
		ID:   check.ID,
		Name: check.Name,
	}
	reg.TCP = check.Address
	reg.Interval = check.Interval.String()
	reg.Timeout = check.Timeout.String()
	return mas.client.Agent().CheckRegister(reg)
}

// Checks returns dict of registred checks
func (nas *NotesAppService) Checks() (map[string]*consul.AgentCheck, error) {
	return nas.client.Agent().Checks()
}

// Services returns dict of registred services
func (nas *NotesAppService) Services() (map[string]*consul.AgentService, error) {
	return nas.client.Agent().Services()
}
