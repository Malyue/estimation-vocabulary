package model

import (
	_config "estimation-vocabulary/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var db *gorm.DB

func init() {
	GetDB()
}

func InitMysql(config *_config.Config) {
	user := config.Mysql.Username
	password := config.Mysql.Password
	address := config.Mysql.Addr
	dbName := config.Mysql.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, address, dbName)
	var err error
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容）
		logger.Config{
			SlowThreshold:             time.Second,   // 慢 SQL 阈值
			LogLevel:                  logger.Silent, // 日志级别
			IgnoreRecordNotFoundError: false,         // 不忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,          // 彩色打印
		},
	)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Println("连接数据库失败")
	}
	fmt.Println("数据库链接成功")

	err = db.AutoMigrate(&Vocabulary{}, &User{})

}

func GetDB() *gorm.DB {
	config := _config.GetConfig()
	if db == nil {
		InitMysql(config)
	}
	return db
}
