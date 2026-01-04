package config

// MysqlConfig MySQL数据库配置
type MysqlConfig struct {
	Path         string
	Port         string
	Config       string
	Dbname       string
	Username     string
	Password     string
	MaxIdleConns int
	MaxOpenConns int
	LogMode      string
	LogZap       bool
	MaxLifetime  int
	AutoCreate   bool
}
