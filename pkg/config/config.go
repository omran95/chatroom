package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Room          *RoomConfig          `mapstructure:"room"`
	Subscriber    *SubscriberConfig    `mapstructure:"subscriber"`
	Cassandra     *CassandraConfig     `mapstructure:"cassandra"`
	Redis         *RedisConfig         `mapstructure:"redis"`
	Kafka         *KafkaConfig         `mapstructure:"kafka"`
	Observability *ObservabilityConfig `mapstructure:"observability"`
}

type RoomConfig struct {
	Http struct {
		Server struct {
			Port    string
			MaxConn int64
		}
	}
	MessageSubscriber struct {
		Topic string
	}
	Grpc struct {
		Client struct {
			SubscriberEndpoint string
		}
	}
}

type SubscriberConfig struct {
	Grpc struct {
		Server struct {
			Port string
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

type RedisConfig struct {
	Password                string
	Addrs                   string
	ExpirationHour          int64
	MinIdleConn             int
	PoolSize                int
	ReadTimeoutMilliSecond  int64
	WriteTimeoutMilliSecond int64
}

type KafkaConfig struct {
	Addrs   string
	Version string
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
	viper.SetDefault("room.messageSubscriber.topic", "room.msg.subscriber"+os.Getenv("HOSTNAME"))
	viper.SetDefault("room.grpc.client.subscriberEndpoint", "localhost:5000")

	viper.SetDefault("subscriber.grpc.server.port", "5000")

	viper.SetDefault("cassandra.hosts", "localhost")
	viper.SetDefault("cassandra.port", 9042)
	viper.SetDefault("cassandra.user", "billy")
	viper.SetDefault("cassandra.password", "p@ssword")
	viper.SetDefault("cassandra.keyspace", "chatroom")

	viper.SetDefault("redis.password", "redis_cluster_password")
	viper.SetDefault("redis.addrs", "localhost:6379,localhost:6380,localhost:6381,localhost:6382,localhost:6383")
	viper.SetDefault("redis.expirationHour", 24)
	viper.SetDefault("redis.minIdleConn", 16)
	viper.SetDefault("redis.poolSize", 64)
	viper.SetDefault("redis.readTimeoutMilliSecond", 3000)
	viper.SetDefault("redis.writeTimeoutMilliSecond", 3000)

	viper.SetDefault("kafka.addrs", "localhost:9092")
	viper.SetDefault("kafka.version", "1.0.0")

	viper.SetDefault("observability.prometheus.port", "")
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
