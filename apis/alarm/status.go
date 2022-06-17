package alarm

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	// "github.com/tidwall/gjson"
)

type Holiday struct {
	Date    string `json:"date"  gorm:"type:varchar(12);"` // 日期,yyyymmdd
	Holiday int    `json:"isHoliday" gorm:"type:int(1);"`  // 节假日，0：上班；1：节假日
}

// 节假日判断
func isHoliday(d string, dbconn *gorm.DB) bool {
	// 节假日表名
	tn := "base_holiday"

	var h Holiday

	// var holidayMap map[string]bool

	// 获取节假日表
	dbconn.Table(tn).Where("date = ?", d).First(&h)

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
func isSchedule(cmd *AlarmCommand, dbconn *gorm.DB) bool {
	// 判断今天是不是节日
	ih := isHoliday(time.Now().Format("20060102"), dbconn)
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
func getURLbyID(id int64, dbconn *gorm.DB) (string, error) {
	// 查询接口表名
	tn := "alarm_apis"

	var a AlarmAPI
	if err := dbconn.Table(tn).Where("id = ?", id).First(&a).Error; err != nil {
		fmt.Println(err)
		return a.APIUrl, err
	}
	return a.APIUrl, nil
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
func getItems(dbconn *gorm.DB) ([]AlarmItem, error) {
	var items []AlarmItem
	if err := dbconn.Table("alarm_items").Where("valid != 'false'").Where("`deleted_at` IS NULL").Preload("Wiki").Preload("Commands").Find(&items).Error; err != nil {
		return items, err
	}
	return items, nil
}
