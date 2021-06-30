package internal

import (
	"fmt"
	"github/rossi1/go-api-microservice-example/domain/entity"
	"net/url"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgreSQL(cfg Config) (*gorm.DB, error) {
	databaseHost := cfg.DB_HOST
	databasePort := cfg.DB_PORT
	databaseUsername := cfg.DB_USER
	databasePassword := cfg.DB_PASSWORD
	databaseName := cfg.DB_NAME
	databaseSSLMode := cfg.DB_SSL_MODE
	// XXX: -

	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(databaseUsername, databasePassword),
		Host:   fmt.Sprintf("%s:%s", databaseHost, databasePort),
		Path:   databaseName,
	}

	q := dsn.Query()
	q.Add("sslmode", databaseSSLMode)

	dsn.RawQuery = q.Encode()

	db, err := gorm.Open(postgres.Open(dsn.String()), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(entity.Category{}, entity.Product{})

	return db, nil

}
