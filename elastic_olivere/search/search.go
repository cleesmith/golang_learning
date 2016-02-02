package search

import (
	// "encoding/json"
	"fmt"
	// "reflect"
	"time"

	"gopkg.in/olivere/elastic.v3" // see: https://github.com/olivere/elastic
)

const (
	TsLayout = "2006-01-02T15:04:05.000Z"
)

type EventU2Record struct {
	RecordType       string    `json:"record_type"`
	Timestamp        time.Time `json:"@timestamp"`
	SensorId         int64     `json:"sensor_id"`
	EventId          int64     `json:"event_id"`
	EventSecond      int64     `json:"event_second"`
	EventMicrosecond int64     `json:"event_microsecond"`
	GeneratorId      int64     `json:"generator_id"`
	SignatureId      int64     `json:"signature_id"`
	SrcIP            string    `json:"src_ip"`
	SPort            int64     `json:"sport"`
	DstIP            string    `json:"dst_ip"`
	DPort            int64     `json:"dport"`
	Protocol         int64     `json:"protocol"`
	Signature        string    `json:"signature"`
}

type PacketU2Record struct {
	RecordType        string    `json:"record_type"`
	Timestamp         time.Time `json:"@timestamp"`
	SensorId          int64     `json:"sensor_id"`
	EventId           int64     `json:"event_id"`
	EventSecond       int64     `json:"event_second"`
	PacketSecond      int64     `json:"packet_second"` // is same as EventSecond?
	PacketMicrosecond int64     `json:"packet_microsecond"`
	IpSrcIP           string    `json:"ip_src_ip"`
	IpDstIP           string    `json:"ip_dst_ip"`
	IpProtocol        int64     `json:"ip_protocol"`
	PacketDump        string    `json:"packet_dump"`
	SrcIP             string
	DstIP             string
	Protocol          int64
}

type ExtraDataU2Record struct {
	RecordType  string    `json:"record_type"`
	Timestamp   time.Time `json:"@timestamp"`
	SensorId    int64     `json:"sensor_id"`
	EventId     int64     `json:"event_id"`
	EventSecond int64     `json:"event_second"`
	EventType   int64     `json:"event_type"`
	EventLength int64     `json:"event_length"`
	XType       int64     `json:"extradata_type"`
	XDataType   int64     `json:"extradata_data_type"`
	XDataLength int64     `json:"extradata_data_length"`
	XData       string    `json:extradata_data"`
}

type AnyU2Record struct {
	Label             string    `json:"label"`
	RecordType        string    `json:"record_type"`
	Timestamp         time.Time `json:"@timestamp"`
	SensorId          int64     `json:"sensor_id"`
	EventId           int64     `json:"event_id"`
	EventSecond       int64     `json:"event_second"`
	EventMicrosecond  int64     `json:"event_microsecond"`
	GeneratorId       int64     `json:"generator_id"`
	SignatureId       int64     `json:"signature_id"`
	SrcIP             string    `json:"src_ip"`
	SPort             int64     `json:"sport"`
	DstIP             string    `json:"dst_ip"`
	DPort             int64     `json:"dport"`
	Protocol          int64     `json:"protocol"`
	Signature         string    `json:"signature"`
	PacketSecond      int64     `json:"packet_second"` // is same as EventSecond?
	PacketMicrosecond int64     `json:"packet_microsecond"`
	PacketDump        string    `json:"packet_dump"`
	EventType         int64     `json:"event_type"`
	EventLength       int64     `json:"event_length"`
	XType             int64     `json:"extradata_type"`
	XDataType         int64     `json:"extradata_data_type"`
	XDataLength       int64     `json:"extradata_data_length"`
	XData             string    `json:extradata_data"`
}

func EsQuery(client *elastic.Client, maxSearchResults int, index string, q string) (*elastic.SearchResult, error) {
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
		// Pretty(true).
		Do()
	return searchResult, err
}

func EsTermsAgg(client *elastic.Client, maxAggResults int, index string, q string, term string) (*elastic.SearchResult, error) {
	// // search like "query": {"match_all":{}} in Sense:
	// // searchResult, err := client.Search().
	// //  Index("unifiedbeat-*").
	// //  Size(MaxSearchResults).
	// //  Do()
	all := elastic.NewMatchAllQuery()
	termAgg := elastic.NewTermsAggregation().Field(term).Size(maxAggResults).OrderByCountDesc()
	aQuery := client.Search().Index(index).Query(all).Size(0).Pretty(true)
	aQuery = aQuery.Aggregation(term, termAgg)
	searchResult, err := aQuery.Do()
	if err != nil {
		fmt.Printf("search error: %#v\n", err)
	}
	return searchResult, err
	// fmt.Printf("\nsearchResult=%T :\n%#v\n", searchResult, searchResult)

	// agg := searchResult.Aggregations
	// if agg == nil {
	// 	fmt.Printf("expected Aggregations != nil; got: nil")
	// }
	// termsAggResult, found := agg.Terms("src_ip")
	// if !found {
	// 	fmt.Printf("expected %v; got: %v", true, found)
	// }
	// // fmt.Printf("\ntermsAggResult=%T :\n%#v\n", termsAggResult, termsAggResult)
	// // fmt.Printf("\nbuckets=%T :\n%#v\n", termsAggResult.Buckets, termsAggResult.Buckets)
	// "termsAggResult.Buckets" is "[]*elastic.AggregationBucketKeyItem"
	// for _, bucket := range termsAggResult.Buckets {
	// 	fmt.Printf("%v \t %v\n", bucket.DocCount, bucket.Key)
	// }

}
