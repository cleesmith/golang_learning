package main

import (
	"encoding/json"
	"fmt"
	"time"
	// "reflect"

	"gopkg.in/olivere/elastic.v3" // see: https://github.com/olivere/elastic
)

const (
	TsLayout         = "2006-01-02T15:04:05.000Z"
	MaxSearchResults = 5
)

type EventU2Record struct {
	Label            string    `json:label`
	Timestamp        time.Time `json:"@timestamp"`
	EventMicrosecond int64     `json:"event_microsecond"`
	EventId          int64     `json:"event_id"`
	GeneratorId      int64     `json:"generator_id"`
	SignatureId      int64     `json:"signature_id"`
	SrcIP            string    `json:"src_ip,omitempty"`
	SPort            int64     `json:"sport,omitempty"`
	DstIP            string    `json:"dst_ip,omitempty"`
	DPort            int64     `json:"dport,omitempty"`
	Protocol         int64     `json:"protocol,omitempty"`
	Signature        string    `json:"signature,omitempty"`
}

type PacketU2Record struct {
	Label             string    `json:label`
	Timestamp         time.Time `json:"@timestamp"`
	PacketMicrosecond int64     `json:"packet_microsecond"`
	EventId           int64     `json:"event_id"`
	PacketSecond      int64     `json:"packet_second,omitempty"`
	PacketDump        string    `json:"packet_dump,omitempty"`
}

func main() {
	client, err := elastic.NewClient(elastic.SetURL("http://192.168.0.31:9200"))
	// client, err := elastic.NewClient()
	if err != nil {
		panic(err)
	}
	fmt.Printf("client=%T=%#v\n", client, client)

	// note: this is really pinging a node so the URL string must be provided:
	info, code, err := client.Ping("http://192.168.0.31:9200").Do()
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("\nPing: Elasticsearch: return code: %d \t version: %s\n", code, info.Version.Number)

	// search via string query: lucene AND/OR/NOT/etc.
	// stringQuery := elastic.NewQueryStringQuery(`packet_payload:localhost AND event_id:2`)
	// stringQuery := elastic.NewQueryStringQuery(`192.168.1.102 AND ftp AND ACK=true`)
	// stringQuery := elastic.NewQueryStringQuery(`"5f 25 25 7c"`)
	stringQuery := elastic.NewQueryStringQuery(`event_id:200`)
	// stringQuery := elastic.NewQueryStringQuery(`*`) // may NOT be empty or blank !!!
	searchResult, err := client.Search().
		Index("unifiedbeat-*").
		Query(stringQuery).
		Size(MaxSearchResults).
		// Pretty(true).
		Do()

	// search via term query
	// termQuery := elastic.NewTermQuery("packet_dump", "localhost")
	// termQuery := elastic.NewTermQuery("packet_payload", "localhost")
	// searchResult, err := client.Search().
	// 	Index("unifiedbeat-*").
	// 	Query(termQuery).
	// 	Size(MaxSearchResults).
	// 	Do()

	// search like "query": {"match_all":{}} in Sense:
	// searchResult, err := client.Search().
	// 	Index("unifiedbeat-*").
	// 	Size(MaxSearchResults).
	// 	Do()
	if err != nil {
		e, ok := err.(*elastic.Error)
		if !ok {
			fmt.Printf("Search failed with unknown error=%T=%#v\n", err, err)
		}
		fmt.Printf("Search failed with status %v and error:\n%s\n", e.Status, e.Details)
		// fmt.Printf("Search failed=%T=%#v\n", e.Details, e.Details)
		return
	}
	fmt.Printf("TotalHits=%T=%#v\n", searchResult.TotalHits(), searchResult.TotalHits())
	fmt.Printf("\nsearchResult=%T :\n%#v\n", searchResult, searchResult)
	fmt.Printf("\nsearchResult.Hits=%T \t count=%#v :\n%#v\n", searchResult.Hits, len(searchResult.Hits.Hits), searchResult.Hits)
	fmt.Printf("\nFound a total of %d unified2 records\n", searchResult.TotalHits())

	// var u2r U2Record

	// note: .Each does NOT give full control over iterating the hits
	// for _, item := range searchResult.Each(reflect.TypeOf(u2r)) {
	// 	fmt.Printf("item=%T=%#v\n", item, item)
	// 	if t, ok := item.(U2Record); ok {
	// 		fmt.Printf("doc=%T=%#v\n", t, t)
	// 	}
	// }

	// iterate through results with full control over each step:
	if searchResult.Hits != nil {
		fmt.Printf("total=%d\n", searchResult.Hits.TotalHits)
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index = "_index"
			// hit.Type = "_type"
			// fmt.Printf("\n-----\nhit: Index=%v Type=%v\n", hit.Index, hit.Type)

			switch hit.Type {
			case "event":
				var event EventU2Record
				err := json.Unmarshal(*hit.Source, &event)
				if err != nil {
					fmt.Printf("Error: deserializing %v type record: %#v!\n", hit.Type, err)
				}
				event.Label = "U2E" // why? future = OSSEC, BRO_LOG, etc.
				fmt.Printf("-----\n%v %v %v event_id=%v gen_sid=%v:%v src_ip=%v sport=%v dst_ip=%v dport=%v proto=%v sig=%v\n",
					event.Label,
					event.Timestamp.UTC().Format(TsLayout),
					event.EventMicrosecond,
					event.EventId,
					event.GeneratorId,
					event.SignatureId,
					event.SrcIP,
					event.SPort,
					event.DstIP,
					event.DPort,
					event.Protocol,
					event.Signature)

			case "packet":
				var packet PacketU2Record
				err := json.Unmarshal(*hit.Source, &packet)
				if err != nil {
					fmt.Printf("Error: deserializing %v type record: %#v!\n", hit.Type, err)
				}
				packet.Label = "U2P" // why? future = OSSEC, BRO_LOG, etc.
				fmt.Printf("-----\n%v %v %v event_id=%v \npacket_dump:\n%s\n",
					packet.Label,
					packet.Timestamp.UTC().Format(TsLayout),
					packet.PacketMicrosecond,
					packet.EventId,
					packet.PacketDump)
			case "extradata":
				fmt.Printf("hit.Type=%T=%#v\n", hit.Type, hit.Type)
				// event.Label = "U2X" // why? future = OSSEC, BRO_LOG, etc.
			default:
				fmt.Printf("hit.Type=%T=%#v\n", hit.Type, hit.Type)
			}

			// deserialize hit.Source into a struct (may also be a map[string]interface{})
			// var t Tweet
			// err := json.Unmarshal(*hit.Source, &t)
			// if err != nil {
			// 	// deserialize failed
			// }
			// fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
		}

		fmt.Printf("TotalHits=%T=%#v\n", searchResult.TotalHits(), searchResult.TotalHits())
		fmt.Printf("\nsearchResult=%T :\n%#v\n", searchResult, searchResult)
		fmt.Printf("\nsearchResult.Hits=%T \t count=%#v :\n%#v\n", searchResult.Hits, len(searchResult.Hits.Hits), searchResult.Hits)
		fmt.Printf("\nFound a total of %d unified2 records\n", searchResult.TotalHits())
	} else {
		fmt.Printf("no hits!\n")
	}

	fmt.Println("\nCount Source IPs:")
	all := elastic.NewMatchAllQuery()
	srcIpAgg := elastic.NewTermsAggregation().Field("src_ip").Size(100).OrderByCountDesc()
	aQuery := client.Search().Index("unifiedbeat-*").Query(all).Size(0).Pretty(true)
	aQuery = aQuery.Aggregation("src_ip", srcIpAgg)
	searchResult, err = aQuery.Do()
	if err != nil {
		fmt.Printf("search error: %#v\n", err)
	}
	// fmt.Printf("\nsearchResult=%T :\n%#v\n", searchResult, searchResult)

	agg := searchResult.Aggregations
	if agg == nil {
		fmt.Printf("expected Aggregations != nil; got: nil")
	}
	termsAggResult, found := agg.Terms("src_ip")
	if !found {
		fmt.Printf("expected %v; got: %v", true, found)
	}
	// fmt.Printf("\ntermsAggResult=%T :\n%#v\n", termsAggResult, termsAggResult)
	// fmt.Printf("\nbuckets=%T :\n%#v\n", termsAggResult.Buckets, termsAggResult.Buckets)
	for _, bucket := range termsAggResult.Buckets {
		fmt.Printf("%v \t %v\n", bucket.DocCount, bucket.Key)
	}

}
