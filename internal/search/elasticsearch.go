package search

import (
	"github.com/olivere/elastic"
)

type ElasticSearcher struct {
	client *elastic.Client
}

func NewElasticSearcher(client *elastic.Client) *ElasticSearcher {

	return &ElasticSearcher{client: client}
}
