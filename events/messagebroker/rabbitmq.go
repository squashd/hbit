package messagebroker

import "github.com/SQUASHD/hbit/config"

func NewRabbitConnectionString(config config.RabbitmqConnnection) string {
	if config.Username == "" {
		config.Username = "guest"
	}
	if config.Password == "" {
		config.Password = "guest"
	}
	if config.Host == "" {
		config.Host = "localhost"
	}
	if config.Port == "" {
		config.Port = "5672"
	}
	return "amqp://" + config.Username + ":" + config.Password + "@" + config.Host + ":" + config.Port
}
