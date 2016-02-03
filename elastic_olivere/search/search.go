package search

import (
	"gopkg.in/olivere/elastic.v3" // see: https://github.com/olivere/elastic
)

func EsQuery(
	client *elastic.Client, maxSearchResults int, index string, q string) (*elastic.SearchResult, error) {

	// search via string query: lucene AND/OR/NOT/etc.
	// stringQuery := elastic.NewQueryStringQuery(`packet_payload:localhost AND event_id:2`)
	// stringQuery := elastic.NewQueryStringQuery(`192.168.1.102 AND ftp AND ACK=true`)
	// stringQuery := elastic.NewQueryStringQuery(`"5f 25 25 7c"`)
	// stringQuery := elastic.NewQueryStringQuery(`event_id:200`)
	// stringQuery := elastic.NewQueryStringQuery(`*`) // may NOT be empty or blank !!!

	stringQuery := elastic.NewQueryStringQuery(q)
	// see: https://github.com/olivere/elastic/blob/release-branch.v3/search.go#L315
	// for "searchResult" struct:
	// type SearchResult struct {
	//  TookInMillis int64         `json:"took"`         // search time in milliseconds
	//  ScrollId     string        `json:"_scroll_id"`   // only used with Scroll and Scan operations
	//  Hits         *SearchHits   `json:"hits"`         // the actual search hits
	//  Suggest      SearchSuggest `json:"suggest"`      // results from suggesters
	//  Aggregations Aggregations  `json:"aggregations"` // results from aggregations
	//  TimedOut     bool          `json:"timed_out"`    // true if the search timed out
	//  Error *ErrorDetails `json:"error,omitempty"` // only used in MultiGet
	// }
	searchResult, err := client.Search().
		Index(index).
		Query(stringQuery).
		Size(maxSearchResults).
		Do()
	return searchResult, err
}

func EsTermsAgg(
	client *elastic.Client, maxAggResults int,
	index string, term string) (*elastic.AggregationBucketKeyItems, error) {

	matchAll := elastic.NewMatchAllQuery()
	termAgg := elastic.NewTermsAggregation().Field(term).Size(maxAggResults).OrderByCountDesc()
	aQuery := client.Search().Index(index).Query(matchAll).Size(0).Pretty(true)
	aQuery = aQuery.Aggregation(term, termAgg)
	searchResult, err := aQuery.Do()
	if err != nil {
		return nil, err
	}
	agg := searchResult.Aggregations
	if agg == nil {
		return nil, err
	}
	termsAggResult, ok := agg.Terms(term)
	if ok {
		return termsAggResult, err
	}

	return nil, err
}
