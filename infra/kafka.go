package infra

import (
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

func KafkaWriter() *kafka.Writer {
	w := &kafka.Writer{
		Addr:  kafka.TCP(KafkaBroker),
		Topic: KafkaTopic,
		Transport: &kafka.Transport{
			SASL: plain.Mechanism{
				Username: KafkaUser,
				Password: KafkaPassword,
			},
		},
	}

	return w
}
