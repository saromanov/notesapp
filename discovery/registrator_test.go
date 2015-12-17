package discovery

import (
   "testing"
   "time"
)

func discoveryApp()(*NotesAppService, error) {
	var err error
	var disc *NotesAppService
	disc, err = CreateDiscovery(&Config{
		  Address: "127.0.0.1:8500",
	})

	if err != nil {
		return nil, err
	}

	return disc, nil
}

func TestRegister(t *testing.T) {
	var err error
	var disc *NotesAppService
	disc, err = discoveryApp()
	if err != nil {
		t.Errorf("%v", err)
	}

	err = disc.Register(&Service{
		   ID: "1234",
		   Name: "Mongo",
		   Address: "127.0.0.1",
		   Port: 27017,
	})

	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestRegisterCheck(t *testing.T) {
	var err error
	var disc *NotesAppService

	disc, err = discoveryApp()
	if err != nil {
		t.Errorf("%v", err)
	}

	err = disc.RegisterCheck(&Check {
		   		ID:   "123",
		   		Name: "mongocheck",
		   		Address: "127.0.0.1:27017",
		   		Interval: 2 * time.Second,
		   		Timeout: 2 * time.Second,
	})

	if err != nil {
		t.Errorf("%v", err)
	}
}