package internal

import (
	"time"

	"github.com/olivere/elastic/v6"
)

func NewElasticSearch(config Config) (*elastic.Client, error) {
	if config.Environment == "local" {
		client, err := elastic.NewClient()
		if err != nil {
			return nil, err
		}

		return client, nil

	}
	client, err := elastic.NewClient(
		elastic.SetURL(config.URL),
		elastic.SetSniff(false),
		elastic.SetBasicAuth(config.USERNAME, config.PASSWORD),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetGzip(true),
	)
	if err != nil {
		return nil, err
	}

	return client, nil

}
