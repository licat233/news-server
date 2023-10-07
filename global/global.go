package global

import (
	"fmt"
	"log"
	"news/dao/query"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func connMysql(dsn string) *query.Query {
	// 配置数据库连接
	// dsn := "root:@tcp(127.0.0.1:3306)/news?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 关闭 GORM 日志
	})
	if err != nil {
		panic("failed to connect database")
	}

	q := query.Use(db)

	// fmt.Println("MySQL successfully connected")

	return q
}

var (
	DB *query.Query
	V  *viper.Viper
)

func init() {
	V = viper.New()
	V.SetConfigName("config") // name of config file (without extension)
	V.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	V.AddConfigPath(".")      // optionally look for config in the working directory
	err := V.ReadInConfig()   // Find and read the config file
	if err != nil {           // Handle errors reading the config file
		log.Fatalln(fmt.Errorf("fatal error config file: %w", err))
	}
	DB = connMysql(V.GetString("datasource"))
}
