job "notesapp" {
	datacenters = ["dc1"]
	type = "service"
	constraint {
		attribute = "$attr.kernel.name"
		value = "linux"
	}
	update {
		stagger = "10s"

		max_parallel = 1
	}

	group "database" {
		task "mongo" {
            driver = "mongo"
            config {
                image = "tobilg/mongodb-marathon"
            }

            resources {
                cpu = 500 
        	    memory = 256

        	    network {
        	        mbits = 10
        		    dynamic_ports = ["mongodb"]
        	    }
            }
        }
	}
}