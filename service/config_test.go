package service

import (
  "testing"
)

func TestCheckConfig(t *testing.T) {
	var cfg Config
	cfg = Config {Name: "testname"}
	if CheckConfig(cfg) == nil {
		t.Errorf("Must return error, cause some fields is black")
	}

	cfg = Config {Name: "testname", RabbitAddr: "AAA"}
	if CheckConfig(cfg) == nil {
		t.Errorf("Must return error, cause MongoAddr field is black")
	}

	cfg = Config {Name: "testname", RabbitAddr: "AAA", MongoAddr: "DDD"}
	if CheckConfig(cfg) != nil {
		t.Errorf("Must not return errors")
	}
}