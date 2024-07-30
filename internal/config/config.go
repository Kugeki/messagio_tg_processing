package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Telegram struct {
		APIToken string `env:"API_TOKEN" env-required:"true"`
		ChatID   int64  `env:"CHAT_ID" env-required:"true"`
	} `env-prefix:"TELEGRAM_"`

	Kafka struct {
		ClientID   string   `env:"CLIENT_ID" env-required:"true"`
		BrokerList []string `env:"BROKER_LIST" env-required:"true"`

		Consumer struct {
			Group  string   `env:"GROUP" env-required:"true"`
			Topics []string `env:"TOPICS" env-required:"true"`
		} `env-prefix:"CONSUMER_"`

		Producer struct {
			Topic string `env:"TOPIC" env-required:"true"`
		} `env-prefix:"PRODUCER_"`
	} `env-prefix:"KAFKA_"`
}

func ReadEnvToConfig() (Config, error) {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	return cfg, err
}
