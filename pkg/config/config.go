package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Room          *RoomConfig          `mapstructure:"room"`
	Cassandra     *CassandraConfig     `mapstructure:"cassandra"`
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

type CassandraConfig struct {
	Hosts    string
	Port     int
	User     string
	Password string
	Keyspace string
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

	viper.SetDefault("cassandra.hosts", "localhost")
	viper.SetDefault("cassandra.port", 9042)
	viper.SetDefault("cassandra.user", "billy")
	viper.SetDefault("cassandra.password", "p@ssword")
	viper.SetDefault("cassandra.keyspace", "chatroom")

	viper.SetDefault("observability.prometheus.port", "8080")
	viper.SetDefault("observability.Tracing.URL", "localhost:4318")

}

func NewConfig() (*Config, error) {
	applyDefaultValues()

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
