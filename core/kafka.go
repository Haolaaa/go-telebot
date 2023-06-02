package core

import (
	"net"
	"strconv"
	"telebot_v2/global"
	"telebot_v2/services"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaService struct {
	Writer *kafka.Writer
}

func NewKafkaService() (*KafkaService, error) {
	topic := "video_release"
	conn, err := kafka.Dial("tcp", "localhost:9092")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return nil, err
	}

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return nil, err
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
		return nil, err
	}

	writer := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaService{
		Writer: writer,
	}, nil
}

const (
	ChatID           = -879407699
	KafkaBroker      = "localhost:9092"
	KafkaTopic       = "video_status"
	RequestURL       = "http://47.98.107.101:9092"
	KafkaReadTimeout = 10 * time.Second
)

func (ks *KafkaService) Reader() {
	bot := global.Bot

	bot.Handle("/run_telebot", services.StartHandler(bot, ks.CreateKafkaReader()))
	bot.Handle("/run_task", services.RunHandler(bot))
	bot.Handle("/run_48", services.AllVideosHandler(bot))

	bot.Handle("/health", services.HealthHandler(bot))

	bot.Handle("/user", services.UserHandler(bot))
	bot.Handle("/chat", services.ChatHandler(bot))
}

func (ks *KafkaService) CreateKafkaReader() *kafka.Reader {
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
