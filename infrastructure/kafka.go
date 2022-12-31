package infrastructure

import (
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func ConnectKafka() (kafkaConn *kafka.Producer, err error) {

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": os.Getenv("KAFKA_URI")})
	if err != nil {
		panic(err)
	}

	return p, err
}
