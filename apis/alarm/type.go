package alarm

import "time"

// 告警接口表
type AlarmAPI struct {
	ID int64 `json:"id" gorm:"type:int(11);primary_key"` // id

	APIType    string `json:"apiType" gorm:"type:string(20)"`     // api接口类型
	APIName    string `json:"apiName" gorm:"type:string(100)"`    // api接口名称
	APIUrl     string `json:"apiUrl" gorm:"type:string(255)"`     // api接口url连接
	APIFullUrl string `json:"apiFullUrl" gorm:"type:string(255)"` // api接口url连接(完整)
}

// 告警查询语句表
type AlarmCommand struct {
	ID int64 `json:"id" gorm:"type:int(11);primary_key"` // id

	Command   string `json:"command" gorm:"type:string(1000)"` // 执行查询命令
	Holiday   int    `json:"holiday" gorm:"type:int(1)"`       // 节假日，0：所有；1：工作日；2：节假日
	Week      string `json:"week" gorm:"type:string(255)"`     // 星期，0-7：日至六
	StartDate string `json:"startDate" gorm:"type:string(8)"`  // 告警起始日期，yyyymmdd
	EndDate   string `json:"endDate" gorm:"type:string(8)"`    // 告警结束日期，yyyymmdd
	StartTime string `json:"startTime" gorm:"type:string(4)"`  // 告警起始时间，hhmi
	EndTime   string `json:"endTime" gorm:"type:string(4)"`    // 告警结束时间，hhmi

	AlarmItemID int64 `json:"-" gorm:"type:int(11)"` // 配置项表关联
}

// 告警配置项
type AlarmItem struct {
	ID int64 `json:"id" gorm:"type:int(11);primary_key"` // id

	Code          string 			`json:"code"              gorm:"type:varchar(10);"`         	// 告警编号
	Group         string 			`json:"group"             gorm:"type:varchar(255);"`       		// 系统组
	Description   string 			`json:"description"       gorm:"type:varchar(255);"` 			// 告警说明
	Description2  string 			`json:"description2"      gorm:"-" `                			// 告警说明(不带value值，用于分析)
	Tags          string 			`json:"tags"              gorm:"type:varchar(1000);"`       	// 特征值
	Valid         string 			`json:"valid"             gorm:"type:varchar(5);"`         		// 是否生效
	Debug         string 			`json:"debug"             gorm:"type:varchar(5);"`         		// 调试模式，不发送告警通知
	NullError     string 			`json:"nullError"         gorm:"type:varchar(5);"`     			// 是否对空值监控
	Interval      int64  			`json:"interval"          gorm:"type:int(10);"`         		// 查询间隔
	Url           int64  			`json:"url"               gorm:"type:int(10);"`              	// API地址
	AlarmInterval int64  			`json:"alarmInterval"     gorm:"type:int(10);"`    				// 告警间隔(保留未用)
	SmsNumbers    string 			`json:"smsNumbers"        gorm:"type:varchar(255);"`  			// 短信接收人
	XmppClients   string 			`json:"xmppClients"       gorm:"type:varchar(255);"` 			// XMPP接收客户端
	WxChats       string 			`json:"wxChats"           gorm:"type:varchar(1000);"`    		// 企业微信机器人
	TriggerNum    int64  			`json:"triggerNum"        gorm:"type:int(5);"`        			// 触发告警次数

	Wiki          AlarmWiki         `json:"wiki"              gorm:"foreignkey:AlarmItemID"`     	// wiki, 一对一关联
	Commands      []AlarmCommand    `json:"commands"          gorm:"foreignkey:AlarmItemID"` 		// wiki, 一对多关联
}

// 告警日志项
type AlarmLog struct {
	ID int64 `json:"-" gorm:"type:int(11);primary_key"` // id

	Code         string   `json:"code" gorm:"type:varchar(10);"`         // 告警编号
	Group        string   `json:"group" gorm:"type:varchar(255);"`       // 系统组
	Description  string   `json:"description" gorm:"type:varchar(255);"` // 告警说明
	Description2 string   `json:"description2" gorm:"-" `                // 告警说明(不带value值，用于分析)
	Tags         string   `json:"tags" gorm:"-"`                         // 特征值
	Debug        string   `json:"debug" gorm:"type:varchar(5);"`         // 调试模式，不发送告警通知，0：否，1：是
	NullError    string   `json:"nullError" gorm:"type:varchar(5);"`     // 是否对空值监控，0：否，1：是
	Interval     int64    `json:"interval" gorm:"type:int(10);"`         // 查询间隔
	TriggerNum   int64    `json:"triggerNum" gorm:"type:int(5);"`        // 触发告警次数
	SmsNumbers   []string `json:"smsNumbers" gorm:"-"`                   // 短信接收人
	XmppClients  []string `json:"xmppClients" gorm:"-"`                  // XMPP接收客户端
	WxChats      []string `json:"wxChats" gorm:"-"`                      // 企业微信机器人

	PromUrl  string `json:"promUrl" gorm:"type:string(255);"`   // api接口url连接
	QueryCmd string `json:"queryCmd" gorm:"type:string(1000);"` // 执行查询命令
	Metric   string `json:"metric" gorm:"type:string(255);"`    // 查询返回度量
	Result   string `json:"result" gorm:"type:string(255);"`    // 查询结果
	CodeCmd  string `json:"codeCmd" gorm:"type:string(20);"`    // Code+CmdNum
	Time     string `json:"time" gorm:"type:string(20);"`       // 记录时间
	WikiUrl  string `json:"-" gorm:"type:string(255);"`         // wiki链接url
	Contacts string `json:"-" gorm:"type:string(100)"`          // 联系人

	CreatedAt time.Time `json:"-"`
}

// 告警wiki
type AlarmWiki struct {
	ID int64 `json:"-" gorm:"type:int(11);primary_key"` // id

	Level             string `json:"level" gorm:"type:string(100)"`             // 等级
	Contacts          string `json:"contacts" gorm:"type:string(100)"`          // 联系人
	Report            string `json:"report" gorm:"type:string(10)"`             // 是否上报
	CorrelationSystem string `json:"correlationSystem" gorm:"type:string(255)"` // 关联系统
	Solution          string `json:"solution" gorm:"type:string(255)"`          // 处理方法
	YwptSysId         string `json:"ywptSysId" gorm:"type:string(255)"`         // 业务系统ID
	YwptItemId        string `json:"ywptItemId" gorm:"type:string(255)"`        // 监视项ID

	AlarmItemID int64 `json:"-" gorm:"type:int(11)"` // 配置项表关联
}
