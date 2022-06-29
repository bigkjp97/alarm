package query

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

type Influxdb struct {
}

func (i *Influxdb) Query(res *Result, httpUrl string, cmd string) (err error) {
	httpUrl = httpUrl + `query?pretty=true&db=app_mon&q=` + url.QueryEscape(cmd)
	httpresp, _ := http.Get(httpUrl)
	body, _ := ioutil.ReadAll(httpresp.Body)

	dataLen := gjson.Get(string(body), "results.0.series").Array()

	if len(dataLen) > 0 { //如无查询到数据则返回空值
		for col := 0; col < len(dataLen); col++ {
			columsStr := "results.0.series." + strconv.Itoa(col) + ".columns"
			valueStr := "results.0.series." + strconv.Itoa(col) + ".values"
			tags := "results.0.series." + strconv.Itoa(col) + ".tags"
			results := gjson.GetMany(string(body), columsStr, valueStr, tags)
			metrics := results[0].Array()
			values := results[1].Array()
			groupTag := results[2].Map()
			valuePosi := values[0]
			var fieldposi []int
			for m, n := range valuePosi.Array() {
				if !strings.Contains(n.Raw, "\"") {
					fieldposi = append(fieldposi, m)
				}
			}
			for _, rownum := range values { //遍历行
				for i := 0; i < len(fieldposi); i++ { //遍历field
					num := fieldposi[i] //field值所在位置
					item := make(map[string]string)
					var val string
					for v_indx, v := range rownum.Array() { //遍历每一行的字段值
						if v_indx > 0 { //时间字段不要
							for m_indx, m_k := range metrics { //遍历度量名字，添加到标签
								if m_indx == v_indx { //度量名称位置与值位置相同

									if Contains(fieldposi, v_indx) == -1 { //若果是field值则跳过
										item[m_k.Raw] = v.Raw //是标签则添加到标签JSON串
									} else {
										if m_indx == num { //如果当前遍历的field与度量名位置相同
											val = v.Raw
										}
									}
								}
							}
						}
					}
					item[`"metrics"`] = metrics[num].Raw
					if len(groupTag) != 0 {
						for k, v := range groupTag {
							item[`"`+k+`"`] = v.Raw
						}
					}
					re := Map2String(item)
					res.Add(re, val)
				}
			}
		}
	}
	return nil
}
func Map2String(m map[string]string) (result string) {
	list := make([]string, 0)
	for k, v := range m {
		t1 := fmt.Sprintf("%s:%s", k, fmt.Sprint(v))
		list = append(list, t1)
	}
	result = strings.Join(list, ",")
	result = "{" + result + "}"
	return
}
func Contains(slice []int, s int) int {
	for index, value := range slice {
		if value == s {
			return index
		}
	}
	return -1
}
