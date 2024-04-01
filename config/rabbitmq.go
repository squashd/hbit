package config

type RabbitmqConfig struct {
	Username string
	Password string
	Host     string
	Port     string
}

func NewRabbitConnectionString(config RabbitmqConfig) string {
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
