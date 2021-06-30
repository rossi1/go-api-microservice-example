package internal

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/olivere/elastic/v6"
	"github.com/sethvargo/go-envconfig"

	"gorm.io/gorm"
)

type Deps struct {
	Config
	DB     *gorm.DB
	Search *elastic.Client
	//Producer      *kafka.Producer
	ProducerTopic string
}

type Config struct {
	DBConfig
	ElasticConfig
	KaftaConfig
	MonitorConfig
	ServerHost  string `env:"SERVER_HOST,required"`
	ServerPort  string `env:"SERVER_PORT,required"`
	Environment string `env:"ENVIRONMENT,required"`
}

type DBConfig struct {
	DB_HOST     string `env:"DB_HOST,required"`
	DB_PORT     string `env:"DB_PORT,required"`
	DB_NAME     string `env:"DB_NAME,required"`
	DB_USER     string `env:"DB_USER,required"`
	DB_PASSWORD string `env:"DB_PASSWORD,required"`
	DB_SSL_MODE string `env:"SSL_MODE,required"`
}

type ElasticConfig struct {
	URL      string
	USERNAME string
	PASSWORD string
}

type KaftaConfig struct {
	HOST  string
	TOPIC string
}

type MonitorConfig struct {
	DSN         string
	DEBUG       bool
	ENVIRONMENT string
}

func GetConfig(ctx context.Context) (Config, error) {
	var config = Config{}

	// Load env variables from a .env file if present
	err := godotenv.Load(".env")
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return config, err
	}

	err = envconfig.Process(ctx, &config)
	if err != nil {
		return config, fmt.Errorf("could not get configuration: %w", err)
	}

	return config, nil
}

func GetDeps(Config Config, DB *gorm.DB, Elastic interface{}, Kafta interface{}, topic string) Deps {

	deps := Deps{
		Config: Config,
		DB:     DB,
		//Search:        Elastic,
		//Producer:      Kafta,
		//ProducerTopic: topic,
	}
	return deps
}
