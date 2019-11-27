package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/chenxuehui1/api_exporter/collector"
	"io"
	"log"
	"net/http"
)

func main() {
	// ... assume a main handler named `handler`
	http.Handle("/metrics", promhttp.Handler())
	// ... other setup
	http.Handle("/query", collector.InstrumentHandler(http.HandlerFunc(Query)))
	http.Handle("/hello", collector.InstrumentHandler(http.HandlerFunc(Hello)))
	log.Fatal(http.ListenAndServe(":9081", nil))
}


// hello
func Hello(w http.ResponseWriter, r *http.Request)  {
	// 请求计数
	//time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	_,_ = io.WriteString(w, "hello world!")
}


// query
func Query(w http.ResponseWriter, r *http.Request)  {
	//模拟业务查询耗时0~1s
	_,_ = io.WriteString(w, "some results")
}
