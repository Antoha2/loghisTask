package config

const MQURL = "amqp://guest:guest@rabbitmq:5672/"

type LogConfig struct {
	Host string
	Port int
}

func GetLogConfig() *LogConfig {

	return &LogConfig{

		Host: "log",
		Port: 5672, //8183,
	}
}
