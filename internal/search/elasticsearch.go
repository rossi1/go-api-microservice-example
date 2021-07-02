package search

import (
	"context"
	"encoding/json"

	"github.com/olivere/elastic/v6"
)

type ElasticSearcher struct {
	client *elastic.Client
}

func NewElasticSearcher(client *elastic.Client) *ElasticSearcher {
	return &ElasticSearcher{client: client}
}

func (es *ElasticSearcher) Search(ctx context.Context, search_obj string) ([]IndexedResult, error) {

	sourceAttr := []string{"name"}
	src := elastic.NewFetchSourceContext(true).Include(sourceAttr...)
	termQuery := elastic.NewTermQuery("name", search_obj)
	searchResult, err := es.client.Search().
		Index("inventories").
		Query(termQuery).
		Sort("created_at", true).
		FetchSourceContext(src).
		From(0).Size(40).
		Pretty(true).
		Do(ctx)

	if err != nil {
		return nil, err
	}
	var t IndexedResult

	tResult := []IndexedResult{}

	if searchResult.TotalHits() > 0 {
		for _, hit := range searchResult.Hits.Hits {

			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				return nil, err
			}

			tResult = append(tResult, t)
		}
	} else {
		return tResult, nil
	}
	return tResult, nil
}

func (es *ElasticSearcher) AutoComplete(ctx context.Context, search_obj string) ([]string, error) {
	productSuggester := elastic.NewCompletionSuggester("category-suggest").
		Fuzziness(2).
		Text(search_obj).
		Field("name.suggest").
		SkipDuplicates(true)

	searchSource := elastic.NewSearchSource().
		Suggester(productSuggester).
		FetchSource(false).
		TrackScores(false)

	searchResult, err := es.client.Search().
		Index("inventories").
		Type("_doc").
		SearchSource(searchSource).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	ProductSuggest := searchResult.Suggest["category-suggest"]

	results := []string{}

	for _, options := range ProductSuggest {
		for _, option := range options.Options {
			results = append(results, option.Text)
		}

	}
	return results, nil
}
