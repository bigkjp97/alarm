package alarm

import (
	"alarm/apis/query"
	"alarm/pkg/utils"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/patrickmn/go-cache"
	"github.com/tidwall/gjson"
)

type Server struct {
	Cfg       *utils.Config // 配置文件
	Pool      *redis.Pool   // redis连接池
	Cache     *cache.Cache  // 缓存
	Dbconn    *gorm.DB      // 数据库连接
	Status_ch chan Status
	Alarm_ch  chan string
}

func (server *Server) RunCheck() {
	// 获取告警项
	items, err := server.getItems()
	if err != nil {
		fmt.Println("get item error")
	}

	// 遍历告警项
	for _, item := range items {
		// 遍历查询语句
		for _, cmd := range item.Commands {
			// 赋值告警编号
			code_num := item.Code + "-" + fmt.Sprint(cmd.ID)
			fmt.Println("获取预警编号")
			fmt.Println(code_num)
			// 获取告警状态
			status, err := server.getStatus(code_num)
			fmt.Println("获取上次检查时间")
			fmt.Println(status.Last_check_time)
			if err != nil {
				status.Code = code_num
			} else if time.Now().Before(status.Last_check_time.Add(time.Duration(item.Interval) * time.Minute)) {
				return
			}

			fmt.Println("获取告警描述")
			fmt.Println(item.Description)
			// 用于分析不匹配value值
			item.Description2 = item.Description

			// 更新查询时间
			status.Last_check_time = time.Now()
			fmt.Println(status.Last_check_time)

			// 更新状态
			server.setStatus(status, time.Duration(item.Interval*item.TriggerNum)*time.Minute)
			fmt.Println("更新状态成功")

			// 判断是否在告警日程
			fmt.Println("检查是否在告警日程内")
			if server.isSchedule(&cmd) {
				fmt.Println("告警日程内")
				api, _ := server.getURLbyID(item.Url)
				var api_url query.ApiSelector
				// 选取查询的数据库
				switch api.APIType {
				case "promethues":
					fmt.Println("查询Prometheus数据库")
					api_url = &query.Promethues{}
				case "influxdb":
					api_url = &query.Influxdb{}
				}

				res := &query.Result{}
				fmt.Println("查询接口地址", api.APIUrl)
				if err := api_url.Query(res, api.APIUrl, cmd.Command); err != nil {
					fmt.Println(err)
				}

				for _, result := range res.Get() {
					fmt.Println("打印查询结果", result)
					metric, value := result["metric"], result["value"]
					fmt.Println(metric, value)
					// metric, value := result.Metric, result.Value
					if item.NullError == "true" {
						if value == "" {
							fmt.Println("空值告警")
							// 触发告警
							if err := server.runNotice(item, api.APIUrl, fmt.Sprint(cmd.ID), cmd.Command, metric, value); err != nil {
								// log.Error("func: srv.check()->srv.notice() send notice error:%v", err)
								fmt.Println(err)
							}
						}
					} else {
						if len(value) > 0 {
							item_cp := item
							reg, _ := regexp.Compile(`\w*\%\%\w*\%\%\w*`)          // %%匹配格式%%
							reg_keys := reg.FindAllString(item_cp.Description, -1) // 告警说明中的匹配替换
							for _, reg_key := range reg_keys {
								rk := strings.Trim(reg_key, "%")
								if strings.ToLower(rk) == "value" {
									item_cp.Description = strings.Replace(item_cp.Description, reg_key, value, -1)
									item_cp.Description2 = strings.Replace(item_cp.Description2, reg_key, "", -1)
								}

								str := gjson.Get(metric, rk).Value()
								item_cp.Description = strings.Replace(item_cp.Description, reg_key, fmt.Sprintf("%v", str), -1)
								item_cp.Description2 = strings.Replace(item_cp.Description2, reg_key, fmt.Sprintf("%v", str), -1)
							}
							// 触发告警
							if err := server.runNotice(item_cp, api.APIUrl, fmt.Sprint(cmd.ID), cmd.Command, metric, value); err != nil {
								// log.Error("func: srv.check()->srv.notice() send notice error:%v", err)
								fmt.Println(err)
							}
						}
					}
				}

			}
		}

	}

}

// 告警推送
func (server *Server) runNotice(item AlarmItem, url, cmd_num, cmd, metric, result string) error {
	code_cmd := item.Code + "-" + cmd_num
	var status Status
	// 获取告警编号的状态
	if val, found := server.Cache.Get(code_cmd); found {
		status = val.(Status)
	}

	status, err := statusFromRedis(server.Cfg, server.Pool, code_cmd)
	if err != nil {
		return err
	}
	status.Alarm_count++

	alarm_interval := item.Interval * item.TriggerNum
	timenow := time.Now()

	if timenow.After(status.Last_alarm_time.Add(time.Duration(alarm_interval*2) * time.Minute)) {
		// 两倍的（告警时间*触发次数）内未产生告警,判定上次告警已经结束,重置告警时间及计数器
		status.Alarm_count = 1
		status.Send_count = 0
		// log.Debug("func: srv.notice(); code:%s reset last_alarm_time(old:%v, new:%v)", codecmd, status.LastAlarmTime, timenow)
	}

	// if status.AlarmCnt >= item.TriggerNum { // 满足告警触发次数
	if true {
		status.Send_count++
		msg := &Msg{
			AlarmLog: AlarmLog{
				Code:         item.Code,
				Group:        item.Group,
				Description:  item.Description,
				Description2: item.Description2,
				Tags:         item.Tags,
				NullError:    item.NullError,
				Interval:     item.Interval,
				TriggerNum:   item.TriggerNum,
				SmsNumbers:   strings.Split(item.SmsNumbers, ","),
				XmppClients:  strings.Split(item.XmppClients, ","),
				WxChats:      strings.Split(item.WxChats, ","),
				Debug:        item.Debug,
				PromUrl:      url,
				QueryCmd:     cmd,
				Metric:       metric,
				Result:       result,
				CodeCmd:      code_cmd,
				Time:         time.Now().Format("2006-01-02 15:04:05"),
				WikiUrl:      item.Wiki.Solution,
				Contacts:     item.Wiki.Contacts,
			},
			Wiki: item.Wiki,
		}

		jstr, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		// log.Debug("func: srv.notice(); send notice: %v, send count:%d", jstr, status.Send_count)
		var alarmCh chan string
		if msg.Debug == "false" {
			alarmCh <- string(jstr)
		}

		fmt.Println(string(jstr))
		// if err := alarm2DB(msg.AlarmLog); err != nil {
		// 	// log.Error("func: srv.notice()-> srv.alarm2DB; send alram to db error:%v", err)
		// }
	}
	// else {
	// 	// log.Debug("func: srv.notice(); code:%s not reach(Interval:%d TriggerNum:%d), current: %d", codecmd, item.Interval, item.TriggerNum, status.AlarmCnt)
	// }
	// status.LastAlarmTime = time.Now()
	// setStatus(status, time.Duration(alarmInterval*2)*time.Minute)
	return nil
}
