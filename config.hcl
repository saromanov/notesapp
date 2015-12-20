Name = "notesapp"
ServicegenChecks = true
MessageBus {
	Addr = "amqp://guest:guest@localhost:5672/"
	Exchange = "notesapp"
}

Services {
	Name = "InsertNote"
	MongoAddr="localhost:27017"
	MongoDBName="notesapp"
	RabbitAddr="amqp://guest:guest@localhost:5672/"
	RabbitExchange="notesapp"
	ServerAddr="localhost:8080"
}

Services {
	Name = "NoteList"
	MongoAddr = "localhost:27017"
	MongoDBName = "notesapp"
	RabbitAddr = "amqp://guest:guest@localhost:5672/"
	RabbitExchange = "notesapp"
	ServerAddr = "localhost:8082"
}

Consul {
	 Address = "127.0.0.1:8500"
}

ConsulChecks {
	ID = "1234"
	Name = "Mongo"
	Address = "127.0.0.1"
	Port = 27017
	Checks {
		 ID = "123"
		 Name = "mongocheck"
		 Address = "127.0.0.1:27017"
	}
}
