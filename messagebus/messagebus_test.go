package messagebus

import(
   "testing"
   "time"
   "fmt"
)

func TestCreateMessageBus(t *testing.T) {
	_, err := CreateMessageBus(&Config{
		Addr: "amqp://guest:guest@localhost:5672/",
	})

	if err != nil {
		t.Errorf(fmt.Sprintf("%v", err))
	}
}

func TestStart(t *testing.T) {
	bus, err := CreateMessageBus(&Config{
		Addr: "amqp://guest:guest@localhost:5672/",
	})

	if err != nil {
		t.Errorf(fmt.Sprintf("%v", err))
	}

	timeChan := time.NewTimer(3 * time.Second).C
    doneChan := make(chan bool)
    go func() {
        bus.Start()
        doneChan <- true
    }()
    
    for {
        select {
        case <- timeChan:
            return
        case <- doneChan:
            t.Errorf("Something went wrong")
      }
    }
}