package es

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"time"

	"github.com/olivere/elastic/v7"

	"github.com/wangzewang/esman/config"
)

type (
	LogMessage struct {
		Message interface{} `json:"log"`
	}
)

var Client *elastic.Client

func Init() {
	config := config.GetConfig()

	var err error
	Client, err = elastic.NewClient(
		elastic.SetURL(config.GetString("es.host")),
		elastic.SetSniff(config.GetBool("es.sniff")),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetMaxRetries(5),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)))

	if err != nil {
		panic(err)
	}
	info, code, err := Client.Ping(config.GetString("es.host")).Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	esversion, err := Client.ElasticsearchVersion(config.GetString("es.host"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)
}

func NewEsStreamQuery(taskName string, stop *int64, resp chan<- string) {

	//task.Status = &Running

queryLoop:
	for {
		//startTimestamp := time.Now().AddDate(-1, 0, 0).UnixNano() / 1e6
		//endTimestamp := time.Now().UnixNano() / 1e6
		//bQ := elastic.NewBoolQuery().Must(elastic.NewQueryStringQuery("*"), rQ)
		//bQ := elastic.NewMatchQuery("kubernetes_namespace.keyword", "kube-system")
		bQ := elastic.NewMatchAllQuery()
		numLine := 10000
		res, err := Client.Search().Index("logstash-2021.02.02").Query(bQ).Size(numLine).SortBy(elastic.NewFieldSort("@timestamp").Asc()).Do(context.Background())
		if err != nil {
			panic(err)
		}

		var lmsg LogMessage

		for _, item := range EachTime(res, reflect.TypeOf(lmsg)) {
			t := item.(LogMessage)
			tempstr := fmt.Sprint(t.Message)
			if strings.Contains(tempstr, "[91m") && strings.Contains(tempstr, "[0m") {
				continue
			}
			if len(tempstr) <= 3 {
				tempstr = ""
			} else {
				var begin = strings.Index(tempstr, "{")
				var end = strings.LastIndex(tempstr, "}")
				tempstr = subString(tempstr, begin+1, end)
			}

			resp <- tempstr
		}

		if atomic.LoadInt64(stop) == 1 {
			break queryLoop
		}
	}
}

func NewEsQuery(taskName string) []string {

	//task.Status = &Running
	var res []string

	//startTimestamp := time.Now().AddDate(-1, 0, 0).UnixNano() / 1e6
	//endTimestamp := time.Now().UnixNano() / 1e6
	//bQ := elastic.NewBoolQuery().Must(elastic.NewQueryStringQuery("*"), rQ)
	//bQ := elastic.NewMatchQuery("kubernetes_namespace.keyword", "kube-system")
	bQ := elastic.NewMatchAllQuery()
	numLine := 10000
	searchRes, err := Client.Search().Index("logstash-2021.02.02").Query(bQ).Size(numLine).SortBy(elastic.NewFieldSort("@timestamp").Asc()).Do(context.Background())
	if err != nil {
		panic(err)
	}

	var lmsg LogMessage

	for _, item := range EachTime(searchRes, reflect.TypeOf(lmsg)) {
		t := item.(LogMessage)
		tempstr := fmt.Sprint(t.Message)
		if strings.Contains(tempstr, "[91m") && strings.Contains(tempstr, "[0m") {
			continue
		}
		if len(tempstr) <= 3 {
			tempstr = ""
		} else {
			var begin = strings.Index(tempstr, "{")
			var end = strings.LastIndex(tempstr, "}")
			tempstr = subString(tempstr, begin+1, end)
		}

		res = append(res, tempstr)
	}

	return res
}
func EachTime(r *elastic.SearchResult, typ reflect.Type) []interface{} {
	if r.Hits == nil || r.Hits.Hits == nil || len(r.Hits.Hits) == 0 {
		return nil
	}
	var slice []interface{}
	for _, hit := range r.Hits.Hits {
		v := reflect.New(typ).Elem()
		if hit.Source == nil {
			slice = append(slice, v.Interface())
			continue
		}
		if err := json.Unmarshal(hit.Source, v.Addr().Interface()); err == nil {
			log := LogMessage{Message: v.Interface()}
			slice = append(slice, log)
		}
	}
	return slice
}

func subString(str string, begin, end int) (substr string) {
	rs := []rune(str)
	lth := len(rs)
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	if end > lth {
		end = lth
	}

	return string(rs[begin:end])
}
