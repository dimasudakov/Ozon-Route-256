package app

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

type Config struct {
	Server struct {
		GatewayPort string `mapstructure:"gateway-port"`
		GrpcPort    string `mapstructure:"grpc-port"`
	} `mapstructure:"server"`
	Database struct {
		Name     string `mapstructure:"name"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Port     string `mapstructure:"port"`
	} `mapstructure:"database"`
	Kafka struct {
		LogsTopicName string   `mapstructure:"logs-topic-name"`
		Brokers       []string `mapstructure:"brokers"`
	} `mapstructure:"kafka"`
}

func InitConfig() (*Config, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	configPath := filepath.Join(currentDir, "configs", "config.yaml")
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
