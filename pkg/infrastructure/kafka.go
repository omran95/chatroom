package infrastructure

import (
	"time"

	"github.com/IBM/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/components/metrics"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/omran95/chatroom/pkg/common"
	"github.com/omran95/chatroom/pkg/config"
	prom "github.com/prometheus/client_golang/prometheus"
)

var (
	logger = watermill.NewStdLogger(
		false,
		false,
	)
)

func NewKafkaPublisher(config *config.Config) (message.Publisher, error) {
	kafkaPublisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   common.GetServerAddrs(config.Kafka.Addrs),
			Marshaler: kafka.DefaultMarshaler{},
		},
		logger,
	)
	if err != nil {
		return nil, err
	}

	return kafkaPublisher, nil
}

func NewKafkaPublisherWithPartitioning(config *config.Config) (message.Publisher, error) {
	kafkaPublisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers: common.GetServerAddrs(config.Kafka.Addrs),
			Marshaler: kafka.NewWithPartitioningMarshaler(func(topic string, msg *message.Message) (string, error) {
				return msg.Metadata.Get("partition_key"), nil
			}),
		},
		logger,
	)
	if err != nil {
		return nil, err
	}

	return kafkaPublisher, nil
}

func NewKafkaSubscriber(config *config.Config) (message.Subscriber, error) {
	saramaConfig := sarama.NewConfig()
	saramaVersion, err := sarama.ParseKafkaVersion(config.Kafka.Version)
	if err != nil {
		return nil, err
	}
	saramaConfig.Version = saramaVersion
	saramaConfig.Consumer.Fetch.Default = 1024 * 1024
	saramaConfig.Consumer.Offsets.AutoCommit.Enable = true
	saramaConfig.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	kafkaSubscriber, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:       common.GetServerAddrs(config.Kafka.Addrs),
			Unmarshaler:   kafka.DefaultMarshaler{},
			ConsumerGroup: config.Kafka.Subscriber.ConsumerGroup,
			InitializeTopicDetails: &sarama.TopicDetail{
				NumPartitions:     config.Kafka.Subscriber.NumPartitions,
				ReplicationFactor: config.Kafka.Subscriber.ReplicationFactor,
			},
			OverwriteSaramaConfig: saramaConfig,
			OTELEnabled:           true,
		},
		logger,
	)
	if err != nil {
		return nil, err
	}

	return kafkaSubscriber, nil
}

func NewBrokerRouter(name string) (*message.Router, error) {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return nil, err
	}

	metricsBuilder := metrics.NewPrometheusMetricsBuilder(prom.DefaultRegisterer, name, "pubsub")
	metricsBuilder.AddPrometheusRouterMetrics(router)

	router.AddMiddleware(
		middleware.CorrelationID,
		middleware.Timeout(time.Second*15),
		middleware.Recoverer,
	)
	return router, nil
}
