package dto

// DbConfig 数据库配置
type DbConfig struct {
	Driver       string `json:"driver"`
	User         string `json:"user"`
	Password     string `json:"password"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	DbName       string `json:"db_name"`
	Chartset     string `json:"charset"`
	ShowSql      bool   `json:"show_sql"`
	LogLevel     int    `json:"log_level"`
	MaxOpenConns int    `json:"max_open_conns"`
	MaxIdleConns int    `json:"max_idle_conns"`
}
