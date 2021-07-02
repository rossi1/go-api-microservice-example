package search

import "context"

type IndexedResult struct {
	Name string `json:"name"`
}

type Searcher interface {
	Search(ctx context.Context, search_obj string) ([]IndexedResult, error)
	AutoComplete(ctx context.Context, search_obj string) ([]string, error)
}
