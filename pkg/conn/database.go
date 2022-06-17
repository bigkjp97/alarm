package conn

import (
	"alarm/apis/alarm"
	"alarm/pkg/utils"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// 获取mysql连接串
func getMysqlConnString(c *utils.Config) string {
	fmt.Printf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		c.Database.Username,
		c.Database.Password,
		c.Database.Db_host,
		c.Database.Db_port,
		c.Database.Db_name,
		c.Database.Charset)
	return fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		c.Database.Username,
		c.Database.Password,
		c.Database.Db_host,
		c.Database.Db_port,
		c.Database.Db_name,
		c.Database.Charset)
}

// 获取相应连接类型的连接串
func getConnString(t string, c *utils.Config) string {
	switch t {
	case "mysql":
		return getMysqlConnString(c)
	}
	return ""
}

func NewDBConn(c *utils.Config) (*gorm.DB, error) {
	// 连接类型
	dt := c.Database.Db_type

	// 连接串
	dtStr := getConnString(dt, c)

	// 获取连接
	conn, err := gorm.Open(dt, dtStr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func TestTables(c *utils.Config) {
	var items []alarm.AlarmItem
	conn, err := NewDBConn(c)
	if err != nil {
		return
	}

	if err := conn.Table("alarm_items").Where("valid != 'false'").Where("`deleted_at` IS NULL").Preload("Wiki").Preload("Commands").Find(&items).Error; err != nil {
		return
	}

	for i := 0; i < 10; i++ {
		fmt.Println(items[i])
	}
	defer conn.Close()
}
