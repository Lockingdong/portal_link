package config

import "fmt"

type DBConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Name     string `mapstructure:"DB_NAME"`
	Schema   string `mapstructure:"DB_SCHEMA"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Port     string `mapstructure:"DB_PORT"`
	Encrypt  string `mapstructure:"DB_ENCRYPT"`
}

func (c *DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable search_path=%s TimeZone=UTC",
		c.Host,
		c.User,
		c.Password,
		c.Name,
		c.Port,
		c.Schema,
	)
}

var dbConfig *DBConfig

func GetDBConfig() *DBConfig {
	return dbConfig
}
