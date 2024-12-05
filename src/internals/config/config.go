package config

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Host   string
	Port   string
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
}

func LoadConfig() Config {
	var cfg Config
	v := viper.New()
	v.SetEnvPrefix("TASKIFY")
	v.SetDefault("HOST", "localhost")
	v.SetDefault("PORT", "3000")
	v.SetDefault("DBHOST", "localhost")
	v.SetDefault("DBPORT", "5432")
	v.SetDefault("DBUSER", "taskify")
	v.SetDefault("DBPASS", "taskify")
	v.SetDefault("DBNAME", "taskify")

	err := v.Unmarshal(&cfg)
	if err != nil {
		logrus.Fatalln(err)
	}

	return cfg
}

func (c *Config) GetStringDatabaseConnection() string {
	return fmt.Sprintf(`postgres://%s:%s@%s:%s/%s`, c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
}
