package models

import (
	"byte_dance_5th/pkg/ymlconfig"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

var (
	DB     *gorm.DB
	Config = ymlconfig.ConfigInit("ByteDance_DataBase", "databaseConfig")
)

func Init() {
	InitDB()
}

func InitDB() {
	var err error

	viper := Config.Viper
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.dbname"),
		viper.GetString("mysql.charset"),
		viper.GetBool("mysql.parseTime"),
		viper.GetString("mysql.loc"),
	)

	DB, err = gorm.Open(mysql.Open(dsn))

	sqlDB, err := DB.DB()
	if err != nil {
		log.Println(err.Error())
	}

	if err != nil {
		log.Println(err.Error())
	}

	err = DB.AutoMigrate(&UserInfo{}, &Video{}, &Comment{}, &Message{}, &User{})
	if err != nil {
		log.Println(err.Error())
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Println(err.Error())
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(viper.GetInt("mysql.maxIdleConns"))

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(viper.GetInt("mysql.maxOpenConns"))

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
}
