package query

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	// "strings"
	"github.com/tidwall/gjson"
)

type Loki struct {
}

func (l *Loki) Query(res *Result, host string, cmd string) (err error) {
	// http://host:ip/api/v1/query?query=<pql>
	http_url := host + `loki/api/v1/query?query=` + url.QueryEscape(cmd)

	http_resp, err := http.Get(http_url)
	fmt.Println(http_url)
	if err != nil || http_resp.StatusCode != 200 {
		if http_resp != nil {
			fmt.Println(err, "查询失败")
			return fmt.Errorf("HTTP请求错误, httpCode: %v", http_resp.StatusCode)
		}
	} else {
		fmt.Println("查询成功")
		body, _ := ioutil.ReadAll(http_resp.Body)
		if body != nil {
			result_type := gjson.GetBytes(body, "data.resultType").String()
			results := gjson.GetBytes(body, "data.result").Array()
			fmt.Println(result_type)
			fmt.Println(results)
			if len(results) != 0 {
				switch result_type {
				case "matrix":
					for _, result := range results {
						metric := gjson.Get(result.String(), "metric").String()
						values := gjson.Get(result.String(), "values").Array()
						var val []gjson.Result
						for _, v := range values {
							val = v.Array()
						}
						value := val[1].String()

						res.Add(metric, value)
					}
				case "vector":
					for _, result := range results {
						metric := gjson.Get(result.String(), "metric").String()
						val := gjson.Get(result.String(), "value").Array()
						value := val[1].String()
						fmt.Println(metric)
						res.Add(metric, value)
					}
				case "streams":
					for _, result := range results {
						var tmp_value []string
						var tmp_stream []string
						stream := gjson.Get(result.String(), "stream").String()
						vals := gjson.Get(result.String(), "values").Array()
						for _,val := range vals {
							value := val.Array()[1].String()
							tmp_value = append(tmp_value, value)
						}
						tmp_stream = append(tmp_stream, stream)
						res.AddStream(tmp_stream, tmp_value)
					}
				}
			} else {
				res.Add("", "")
			}
			return nil
		}
	}
	return nil
}
