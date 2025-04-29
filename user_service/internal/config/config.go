package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	MongoURI   string
	ServerPort string
	JWTSecret  string
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs") // looks inside configs/ folder

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	config := &Config{
		MongoURI:   viper.GetString("mongo.uri"),
		ServerPort: viper.GetString("server.port"),
		JWTSecret:  viper.GetString("jwt.secret"),
	}

	return config
}
