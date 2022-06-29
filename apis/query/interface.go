package query

import "fmt"

type Result struct {
	results        []map[string]string
	stream_results []map[string][]string
}

func (r *Result) Add(metric, value string) {
	res := make(map[string]string)
	res["metric"], res["value"] = metric, value

	r.results = append(r.results, res)
	fmt.Println(r.results)
}

func (r *Result) AddStream(stream, values []string) {
	res := make(map[string][]string)
	res["stream"], res["values"] = stream, values

	r.stream_results = append(r.stream_results, res)
	fmt.Println(r.stream_results)
}

func (r *Result) Get() (results []map[string]string) {
	return r.results
}

type ApiSelector interface {
	Query(res *Result, url string, cmd string) (err error)
}
