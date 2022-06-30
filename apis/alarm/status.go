package alarm

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	// "github.com/jinzhu/gorm"
	// "github.com/tidwall/gjson"
)

// 告警状态
type Status struct {
	Code            string    // 告警编号
	Last_check_time time.Time `json:"last_check_time"` // 最后查询时间
	Last_alarm_time time.Time `json:"last_alarm_time"` // 最后告警时间
	Alarm_count     int64     `json:"alarm_count"`     // 告警计数
	Send_count      int64     `json:"send_count"`      // 发送计数
}

type Holiday struct {
	Date    string `json:"date"  gorm:"type:varchar(12);"` // 日期,yyyymmdd
	Holiday int    `json:"isHoliday" gorm:"type:int(1);"`  // 节假日，0：上班；1：节假日
}

// 节假日判断
func (server *Server) isHoliday(d string) bool {
	// 节假日表名
	tn := "base_holiday"

	var h Holiday

	// var holidayMap map[string]bool

	// 获取节假日表
	server.Dbconn.Table(tn).Where("date = ?", d).First(&h)

	fmt.Println(h)
	// 在redis设置一个键，当天过期，没这个键就查表
	return h.Holiday == 1
}

type Schedule struct {
	Is_check bool // 节假日是否符合
	Is_week  bool // 星期是否符合
	Is_date  bool // 日期是否符合
	Is_time  bool // 时间是否符合
}

// 节假日、工作日判断流程
func (server *Server) isSchedule(cmd *AlarmCommand) bool {
	// 判断今天是不是节日
	ih := server.isHoliday(time.Now().Format("20060102"))
	// 引用今天是否告警的判断结构
	s := Schedule{}

	switch cmd.Holiday {
	case 0:
		s.Is_check = true
	case 1:
		if !ih {
			s.Is_check = true
		}
	case 2:
		if ih {
			s.Is_check = true
		}
	}

	// 星期判断
	weekdays := strings.Split(cmd.Week, ",")
	for _, w := range weekdays {
		wi, _ := strconv.Atoi(w) // string转换为int
		if wi == int(time.Now().Weekday()) {
			s.Is_week = true
		}
	}

	// 日期判断
	d := time.Now().Format("20060102")
	if cmd.StartDate <= d && d <= cmd.EndDate {
		s.Is_date = true
	}

	// 时间判断
	t := time.Now().Format("1504")
	if cmd.StartTime <= t && t <= cmd.EndTime {
		s.Is_time = true
	}

	return s.Is_check && s.Is_week && s.Is_date && s.Is_time
}

// 获取api连接
func (server *Server) getURLbyID(id int64) (AlarmAPI, error) {
	// 查询接口表名
	tn := "alarm_apis"

	var a AlarmAPI
	if err := server.Dbconn.Table(tn).Where("id = ?", id).First(&a).Error; err != nil {
		fmt.Println(err)
		return a, err
	}
	return a, nil
}

// func getURLs() (err error) {
// 	var apis []AlarmAPI
// 	if err = srv.DbConn.Table("alarm_apis").Find(&apis).Error; err != nil {
// 		return err
// 	}

// 	for _, a := range apis {
// 		mutex.Lock()
// 		srv.urlAPIMap[a.ID] = a
// 		mutex.Unlock()
// 	}
// 	return nil
// }

// 获取所有可用告警配置项
func (server *Server) getItems() ([]AlarmItem, error) {
	var items []AlarmItem
	if err := server.Dbconn.Table("alarm_items_test").Where("valid != 'false'").Where("`deleted_at` IS NULL").Preload("Wiki").Preload("Commands").Find(&items).Error; err != nil {
		return items, err
	}
	return items, nil
}

// 获取状态缓存
func (server *Server) getStatus(k string) (Status, error) {
	if val, found := server.Cache.Get(k); found {
		fmt.Println(found)
		fmt.Println(val.(Status))
		return val.(Status), nil
	}

	return server.statusFromRedis(k)
}

// 更新状态缓存
func (server *Server) setStatus(s Status, exp time.Duration) {
	// 全局的管道
	server.Cache.Set(s.Code, s, exp)
	server.Status_ch <- s
	fmt.Println("更新状态成功")
}
