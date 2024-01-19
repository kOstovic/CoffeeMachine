package config

import (
	"github.com/spf13/viper"
)

var Configuration Config

type Config struct {
	Database struct {
		DB_TYPE             string `mapstructure:"type"`
		DB_HOST             string `mapstructure:"host"`
		DB_PORT             string `mapstructure:"port"`
		DB_USER             string `mapstructure:"user"`
		DB_PASSWORD         string `mapstructure:"password"`
		DB_PARAMETERS       string `mapstructure:"parameters"`
		DBNAME_INGREDIENT   string `mapstructure:"dbname_ingredient"`
		DBNAME_DENOMINATION string `mapstructure:"dbname_denomination"`
		DBNAME_DRINKS       string `mapstructure:"dbname_drinks"`
		INITIALIZED         string `mapstructure:"initialized"`
	} `mapstructure:"database"`
	Log struct {
		LOG_LEVEL string `mapstructure:"level"`
	} `mapstructure:"log"`
	Auth struct {
		USERNAME string `mapstructure:"username"`
		PASSWORD string `mapstructure:"password"`
	} `mapstructure:"auth"`
}

func LoadConfig(config *Config) error {
	viper.AutomaticEnv()
	viper.SetConfigName("config") // Name of the configuration file without extension
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/bin") // Path to look for the configuration file
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return err
	}

	return nil
}
