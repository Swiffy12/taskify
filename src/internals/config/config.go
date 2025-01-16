package config

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Host   string
	Port   string
	DBHost string `mapstructure:"DATABASE_HOST"`
	DBPort string `mapstructure:"DATABASE_PORT"`
	DBUser string `mapstructure:"DATABASE_USER"`
	DBPass string `mapstructure:"DATABASE_PASS"`
	DBName string `mapstructure:"DATABASE_NAME"`
}

func LoadConfig() Config {
	var cfg Config
	v := viper.New()
	v.SetEnvPrefix("TASKIFY")
	v.SetDefault("HOST", "localhost")
	v.SetDefault("PORT", "3000")
	v.SetDefault("DATABASE_HOST", "localhost")
	v.SetDefault("DATABASE_PORT", "5432")
	v.SetDefault("DATABASE_USER", "taskify")
	v.SetDefault("DATABASE_PASS", "taskify")
	v.SetDefault("DATABASE_NAME", "taskify")
	v.SetDefault("JWT_SECRET_KEY", "fghiwemDUELdF01d3")
	v.SetDefault("USER_ID_KEY", "userId")
	v.AutomaticEnv()

	err := v.Unmarshal(&cfg)
	if err != nil {
		logrus.Fatalln(err)
	}

	return cfg
}

func (c *Config) GetStringDatabaseConnection() string {
	return fmt.Sprintf(`postgres://%s:%s@%s:%s/%s`, c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
}
