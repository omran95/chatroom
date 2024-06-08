package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Room          *RoomConfig          `mapstructure:"room"`
	Observability *ObservabilityConfig `mapstructure:"observability"`
}

type RoomConfig struct {
	Http struct {
		Server struct {
			Port    string
			MaxConn int64
		}
	}
}

type ObservabilityConfig struct {
	Prometheus struct {
		Port string
	}
	Tracing struct {
		URL string
	}
}

func applyDefaultValues() {
	viper.SetDefault("room.http.server.port", "3000")
	viper.SetDefault("room.http.server.maxConn", 20000)
	viper.SetDefault("observability.prometheus.port", "8080")
	viper.SetDefault("observability.Tracing.URL", "http://localhost:5050")

}

func NewConfig() (*Config, error) {
	applyDefaultValues()

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
