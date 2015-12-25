package utils

import (
	"gopkg.in/mgo.v2"
	"github.com/streadway/amqp"
	"time"
	"errors"
)

var (
	errMongoDown = errors.New("Service mongodb is not available")
	errRabbitDown = errors.New("Service rabbit is not available")
)

// Helful functions for wating infrastracture

func WaitForMongo(addr string) error {
	doneChan := make(chan bool)
	completeChan := make(chan bool)
	go func() {
		time.Sleep(time.Second * 20)
		doneChan <- true
	}()

	go func() {
		for {
			_, err := mgo.Dial(addr)
			if err == nil {
				completeChan <- true
			}
		}
	}()

	for {
		select {
		case <-doneChan:
			return errMongoDown
		case <-completeChan:
			return nil
		}
	}
}

func WaitForRabbit(addr string) error {
	doneChan := make(chan bool)
	completeChan := make(chan bool)
	go func() {
		time.Sleep(time.Second * 20)
		doneChan <- true
	}()

	go func() {
		for {
			_, err := amqp.Dial(addr)
			if err == nil {
				completeChan <- true
			}
		}
	}()

	for {
		select {
		case <-doneChan:
			return errRabbitDown
		case <-completeChan:
			return nil
		}
	}
}
