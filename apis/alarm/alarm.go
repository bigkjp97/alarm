package alarm

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"alarm/pkg/utils"

	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/patrickmn/go-cache"
)

func RunCheck(c *utils.Config, r *redis.Pool, cache *cache.Cache, dbconn *gorm.DB,sch chan Status) {
	// 获取告警项
	items, err := getItems(dbconn)
	if err != nil {
		fmt.Println("get item error")
	}

	// 遍历告警项
	for _, item := range items {
		// 遍历查询语句
		for _, cmd := range item.Commands {
			// 赋值告警编号
			code_num := item.Code + "-" + fmt.Sprint(cmd.ID)
			// 获取告警状态
			status, err := getStatus(c, r, cache, code_num)
			if err != nil {
				status.Code = code_num
			} else if time.Now().Before(status.Last_check_time.Add(time.Duration(item.Interval) * time.Minute)) {
				return
			}

			// 用于分析不匹配value值
			item.Description2 = item.Description

			// 更新查询时间
			status.Last_check_time = time.Now()

			// 更新状态
			setStatus(cache, status, sch, time.Duration(item.Interval*item.TriggerNum)*time.Minute)

		}

	}

}

// 告警推送
func runNotice(item AlarmItem, cache *cache.Cache, c *utils.Config, r *redis.Pool, url, cmd_num, cmd, metric, result string) error {
	code_cmd := item.Code + "-" + cmd_num
	var status Status
	// 获取告警编号的状态
	if val, found := cache.Get(code_cmd); found {
		status = val.(Status)
	}

	status, err := statusFromRedis(c, r, code_cmd)
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
