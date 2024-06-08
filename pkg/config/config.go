package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Room *RoomConfig `mapstructure:"room"`
}

type RoomConfig struct {
	Http struct {
		Server struct {
			Port    string
			MaxConn int64
		}
	}
}

func applyDefaultValues() {
	viper.SetDefault("room.http.server.port", "5001")
	viper.SetDefault("room.http.server.maxConn", 200)
}

func NewConfig() (*Config, error) {
	applyDefaultValues()

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
