package utils

import (
	"fmt"
	"strconv"
	"tiny-admin-api-serve/entity/dto"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBUtil struct {
	DB *gorm.DB
}

// Db  全局变量, 外部使用utils.Db来访问
var Db DBUtil

func init() {
	viper.SetConfigFile("./config/config.yaml")

	err := viper.ReadInConfig() // 读取配置信息
	if err != nil {
		// 读取配置信息失败
		fmt.Printf("viper.ReadInConfig failed, err:%v\n", err)
		panic(err)
		return
	}
	dbConfig := dto.DbConfig{
		Driver:       viper.GetString("datasource.type"),
		Host:         viper.GetString("datasource.host"),
		Port:         viper.GetInt("datasource.port"),
		User:         viper.GetString("datasource.user"),
		Password:     viper.GetString("datasource.password"),
		DbName:       viper.GetString("datasource.dbname"),
		Chartset:     viper.GetString("datasource.chartset"),
		ShowSql:      viper.GetBool("datasource.show_sql"),
		LogLevel:     viper.GetInt("datasource.log_level"),
		MaxOpenConns: viper.GetInt("datasource.max_open_conns"),
		MaxIdleConns: viper.GetInt("datasource.max_idle_conns")}
	dsn := buildDatabaseDSN(dbConfig)
	dd, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(dbConfig.LogLevel)),
	})
	if err != nil {
		panic("[SetupDefaultDatabase#newConnection error]: " + err.Error() + " " + dsn)
	}

	sqlDB, err := dd.DB()
	if err != nil {
		panic("[SetupDefaultDatabase#newConnection error]: " + err.Error() + " " + dsn)
	}
	if dbConfig.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	}
	if dbConfig.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	}

	//初始化全局DB连接
	Db = DBUtil{DB: dd}

}

func buildDatabaseDSN(config dto.DbConfig) string {
	switch config.Driver {
	case "mysql":
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
			config.User,
			config.Password,
			config.Host,
			strconv.Itoa(config.Port),
			config.DbName,
		)
	case "postgres":
		return fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s password=%s",
			config.Host,
			strconv.Itoa(config.Port),
			config.User,
			config.DbName,
			config.Password,
		)
	case "sqlite3":
		return config.DbName
	case "mssql":
		return fmt.Sprintf(
			"sqlserver://%s:%s@%s:%s?database=%s",
			config.User,
			config.Password,
			config.Host,
			strconv.Itoa(config.Port),
			config.DbName,
		)
	}
	panic("DB Connection not supported:" + config.Driver)
	return ""
}
