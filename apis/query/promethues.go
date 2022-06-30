package query

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/tidwall/gjson"
)

type Promethues struct {
}

// Promethues查询方法
func (p *Promethues) Query(res *Result, httpUrl string, cmd string) (err error) {
	// http://host:ip/api/v1/query?query=<pql>
	httpUrl = httpUrl + `api/v1/query?query=` + url.QueryEscape(cmd)
	httpresp, err := http.Get(httpUrl)
	if err != nil || httpresp.StatusCode != 200 {
		if httpresp != nil {
			return fmt.Errorf("HTTP请求错误, httpCode: %v", httpresp.StatusCode)
		}
	} else {
		body, _ := ioutil.ReadAll(httpresp.Body)
		if body != nil {
			resultType := gjson.GetBytes(body, "data.resultType").String()
			results := gjson.GetBytes(body, "data.result").Array()

			if len(results) != 0 {
				switch resultType {
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

						res.Add(metric, value)
					}
				case "scalar", "string":
					value := results[1].String()

					res.Add("", value)
				}
			} else {
				res.Add("", "")
			}
			return nil
		}
	}
	return nil
}
