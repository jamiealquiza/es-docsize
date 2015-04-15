package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var (
	nodeIp   string
	nodePort string
)

var metrics struct {
	Indices struct {
		Docs struct {
			Count int64 `json:"count"`
		} `json:"docs"`
	} `json:"indices"`
	Nodes struct {
		Fs struct {
			TotalInBytes     int64 `json:"total_in_bytes"`
			AvailableInBytes int64 `json:"available_in_bytes"`
		} `json:"fs"`
	} `json:"nodes"`
}

func init() {
	flag.StringVar(&nodeIp, "ip", "127.0.0.1", "ElasticSearch IP address")
	flag.StringVar(&nodePort, "port", "9200", "ElasticSearch port")
	flag.Parse()
}

func main() {
	resp, err := http.Get("http://" + nodeIp + ":" + nodePort + "/_cluster/stats")
	if err != nil {
		os.Exit(1)
	}

	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		os.Exit(1)
	}

	json.Unmarshal(contents, &metrics)

	stats := make(map[string]interface{})

	stats["storage_used"] = metrics.Nodes.Fs.TotalInBytes - metrics.Nodes.Fs.AvailableInBytes
	stats["docs"] = metrics.Indices.Docs.Count
	stats["doc_size"] = fmt.Sprintf("%.2f", float64(stats["storage_used"].(int64))/float64(stats["docs"].(int64)))

	out, _ := json.Marshal(stats)
	fmt.Println(string(out))
}
