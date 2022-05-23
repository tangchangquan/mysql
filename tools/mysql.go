package tools

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
)

var Db *gorm.DB

type DbConfig struct {
	Username      string
	Password      string
	Host          string
	Port          int
	Prefix        string
	Extend        string
	SingularTable bool
	LogColorful   bool
	LogLevel      LogLevel
	MaxIdleConns  int
	MaxOpenConns  int
	DbName        string
}

type LogLevel = logger2.LogLevel

// InitMysql @Title 初始化Mysql数据库
func InitMysql(dbConfig *DbConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DbName)
	if dbConfig.Extend != "" {
		dsn += "?" + dbConfig.Extend
	}
	mysqlConfig := mysql.Config{
		DSN: dsn, // DSN data source name
	}
	if Db, err = gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   dbConfig.Prefix,         // 表名前缀
			SingularTable: !dbConfig.SingularTable, // 使用单数表名
		},
		Logger: logger2.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger2.Config{
				Colorful: !dbConfig.LogColorful,
				LogLevel: dbConfig.LogLevel, // Log level
			},
		),
	}); err != nil {
		return
	} else {
		sqlDB, _ := Db.DB()
		sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
		sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
		return
	}
}
