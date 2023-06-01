package core

import (
	"net"
	"strconv"
	"telebot_v2/global"
	"telebot_v2/services"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func NewKafka() {
	topic := "video_release"
	conn, err := kafka.Dial("tcp", "localhost:9092")
	if err != nil {
		global.LOG.Error("Kafka dial failed", zap.Error(err))
	}

	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		global.LOG.Error("Kafka get controller failed", zap.Error(err))
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		global.LOG.Error("Kafka dial failed", zap.Error(err))
	}

	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		global.LOG.Error("Kafka create topic failed", zap.Error(err))
	}
}

func Writer() *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Balancer: &kafka.LeastBytes{},
	}
}

const (
	ChatID           = -879407699
	KafkaBroker      = "localhost:9092"
	KafkaTopic       = "video_status"
	RequestURL       = "http://47.98.107.101:9092"
	KafkaReadTimeout = 10 * time.Second
)

func Reader() {
	bot := global.Bot

	bot.Handle("/start", services.StartHandler(bot, CreateKafkaReader()))

	bot.Handle("/user", services.UserHandler(bot))
	bot.Handle("/chat", services.ChatHandler(bot))
	bot.Handle("/health", services.HealthHandler(bot))
	bot.Handle("/all", services.AllVideosHandler(bot))
}

func CreateKafkaReader() *kafka.Reader {
	return kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:   []string{KafkaBroker},
			Topic:     KafkaTopic,
			Partition: 0,
			MinBytes:  10e3,
			MaxBytes:  10e6,
		},
	)
}
