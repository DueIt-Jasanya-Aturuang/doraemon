package helper

import (
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"

	"github.com/DueIt-Jasanya-Aturuang/doraemon/infra/config"
)

func KafkaWriter() *kafka.Writer {
	w := &kafka.Writer{
		Addr:  kafka.TCP(config.KafkaBroker),
		Topic: config.KafkaTopic,
		Transport: &kafka.Transport{
			SASL: plain.Mechanism{
				Username: config.KafkaUser,
				Password: config.KafkaPassword,
			},
		},
	}

	return w
}
