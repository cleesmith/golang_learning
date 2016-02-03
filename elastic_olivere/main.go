package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"./search"
	"gopkg.in/olivere/elastic.v3" // see: https://github.com/olivere/elastic
)

const (
	EsHostPort       = "http://192.168.0.31:9200"
	TsLayout         = "2006-01-02T15:04:05.000Z"
	MaxSearchResults = 5
	TopTen           = 10
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

func main() {
	client, err := elastic.NewClient(elastic.SetURL(EsHostPort))
	if err != nil {
		panic(err)
	}
	fmt.Printf("client=%T=%#v\n", client, client)

	// note: this is really pinging a node so the URL string must be provided:
	info, code, err := client.Ping(EsHostPort).Do()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nPing: Elasticsearch: return code: %d \t version: %s\n", code, info.Version.Number)

	searchResult, err := search.EsQuery(client, MaxSearchResults, "unifiedbeat-*", `event_id:200`)
	if err != nil {
		e, ok := err.(*elastic.Error)
		if !ok {
			fmt.Printf("Search failed with unknown error=%T=%#v\n", err, err)
		}
		fmt.Printf("Search failed with status %v and error:\n%s\n", e.Status, e.Details)
		return
	}
	fmt.Printf("TotalHits=%T=%#v\n", searchResult.TotalHits(), searchResult.TotalHits())
	fmt.Printf("\nsearchResult=%T :\n%#v\n", searchResult, searchResult)
	fmt.Printf("\nsearchResult.Hits=%T \t count=%#v :\n%#v\n", searchResult.Hits, len(searchResult.Hits.Hits), searchResult.Hits)
	fmt.Printf("\nFound a total of %d unified2 records\n", searchResult.TotalHits())

	// iterate through results with full control over each step:
	if searchResult.Hits != nil {
		fmt.Printf("total=%d\n", searchResult.Hits.TotalHits)
		// see: https://github.com/olivere/elastic/blob/release-branch.v3/search.go#L354
		// type SearchHits struct {
		// 	 TotalHits int64        `json:"total"`     // total number of hits found
		// 	 MaxScore  *float64     `json:"max_score"` // maximum score of all hits
		// 	 Hits      []*SearchHit `json:"hits"`      // the actual hits returned
		// }
		// // SearchHit is a single hit.
		// type SearchHit struct {
		// 	 Score          *float64                       `json:"_score"`          // computed score
		// 	 Index          string                         `json:"_index"`          // index name
		// 	 Type           string                         `json:"_type"`           // type meta field
		// 	 Id             string                         `json:"_id"`             // external or internal
		// 	 Uid            string                         `json:"_uid"`            // uid meta field (see MapperService.java for all meta fields)
		// 	 Timestamp      int64                          `json:"_timestamp"`      // timestamp meta field
		// 	 TTL            int64                          `json:"_ttl"`            // ttl meta field
		// 	 Routing        string                         `json:"_routing"`        // routing meta field
		// 	 Parent         string                         `json:"_parent"`         // parent meta field
		// 	 Version        *int64                         `json:"_version"`        // version number, when Version is set to true in SearchService
		// 	 Sort           []interface{}                  `json:"sort"`            // sort information
		// 	 Highlight      SearchHitHighlight             `json:"highlight"`       // highlighter information
		// 	 Source         *json.RawMessage               `json:"_source"`         // stored document source
		// 	 Fields         map[string]interface{}         `json:"fields"`          // returned fields
		// 	 Explanation    *SearchExplanation             `json:"_explanation"`    // explains how the score was computed
		// 	 MatchedQueries []string                       `json:"matched_queries"` // matched queries
		// 	 InnerHits      map[string]*SearchHitInnerHits `json:"inner_hits"`      // inner hits with ES >= 1.5.0
		// }
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index is "_index"
			// hit.Type is "_type"
			switch hit.Type {
			case "event":
				var event EventU2Record
				// note: hit.Source is a "*json.RawMessage" == "type RawMessage []byte"
				// fmt.Printf("\n-----\nhit.Source=%T=%#v\n", hit.Source, hit.Source)
				// raw, err := hit.Source.MarshalJSON()
				// if err != nil {
				// 	return
				// }
				// fmt.Printf("\n-----\nraw=%T=%#v\n", raw, raw)
				// err = json.Unmarshal(raw, &event)
				err := json.Unmarshal(*hit.Source, &event)
				if err != nil {
					fmt.Printf("Error: deserializing %v type record: %#v!\n", hit.Type, err)
				}
				// fmt.Printf("\n-----\nevent=%T=%#v\n", event, event)
				fmt.Printf("-----\n%v %v %v event_id=%v gid_sid=%v:%v src_ip=%v sport=%v dst_ip=%v dport=%v proto=%v sig=%v\n",
					event.RecordType,
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
				packet.SrcIP = packet.IpSrcIP
				packet.DstIP = packet.IpDstIP
				packet.Protocol = packet.IpProtocol
				fmt.Printf("-----\n%v %v %v event_id=%v SrcIP=%v DstIp=%v Protocol=%v\npacket_dump:\n%s\n",
					packet.RecordType,
					packet.Timestamp.UTC().Format(TsLayout),
					packet.PacketMicrosecond,
					packet.EventId,
					packet.SrcIP,
					packet.DstIP,
					packet.Protocol,
					packet.PacketDump)

			case "extradata":
				var extradata ExtraDataU2Record
				err := json.Unmarshal(*hit.Source, &extradata)
				if err != nil {
					fmt.Printf("Error: deserializing %v type record: %#v!\n", hit.Type, err)
				}
				fmt.Printf("-----\n%v %v %v event_id=%v\n",
					extradata.RecordType,
					extradata.Timestamp.UTC().Format(TsLayout),
					extradata.SensorId,
					extradata.EventId,
					extradata.EventSecond,
					extradata.XData)

			default:
				fmt.Printf("Error: unknown record type of '%#v'\n", hit.Type)
			}

		}

		fmt.Printf("TotalHits=%T=%#v\n", searchResult.TotalHits(), searchResult.TotalHits())
		fmt.Printf("\nsearchResult=%T :\n%#v\n", searchResult, searchResult)
		fmt.Printf("\nsearchResult.Hits=%T \t count=%#v :\n%#v\n", searchResult.Hits, len(searchResult.Hits.Hits), searchResult.Hits)
		fmt.Printf("\nFound a total of %d unified2 records\n", searchResult.TotalHits())

		fmt.Println("\n***** AnyU2Record ...")
		var u2r AnyU2Record
		// note: .Each does NOT give full control over iterating the hits
		for i, hit := range searchResult.Each(reflect.TypeOf(u2r)) {
			// why Label? future = OSSEC, BRO_LOG, etc.
			// Label = "U2E" // "U2P", "U2X"
			fmt.Printf("\n%v. hit=%#v\n", i, hit)
			// 	if t, ok := hit.(U2Record); ok {
			// 		fmt.Printf("doc=%T=%#v\n", t, t)
			// 	}
		}
		fmt.Println("... AnyU2Record ********************************************************")
	} else {
		fmt.Printf("no hits!\n")
	}

	// aField := "event_second"
	aField := "signature.raw"
	fmt.Printf("\nCount\t|  %s\n", aField)
	fmt.Println("=====\t|  ==================")
	aggResult, err := search.EsTermsAgg(client, TopTen, "unifiedbeat-*", aField)
	if aggResult != nil {
		// fmt.Printf("\naggResult=%T :\n%#v\n", aggResult, aggResult)
		// fmt.Printf("\nbuckets=%T :\n%#v\n", aggResult.Buckets, aggResult.Buckets)
		for _, bucket := range aggResult.Buckets {
			if aField == "event_second" {
				ut := time.Unix(int64(bucket.Key.(float64)), 0)
				fmt.Printf("%v\t|  %v\n", bucket.DocCount, ut.UTC().Format(TsLayout))
			} else {
				fmt.Printf("%v\t|  %v\n", bucket.DocCount, bucket.Key)
			}
		}
		fmt.Printf("#buckets=%v\n", len(aggResult.Buckets))
	} else {
		fmt.Printf("\nEmpty: aggResult=%T=%#v\n", aggResult, aggResult)
	}

}
