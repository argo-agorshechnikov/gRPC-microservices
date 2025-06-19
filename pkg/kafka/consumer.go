package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(topic, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: GetBrokers(),
			GroupID: groupID,
			Topic:   topic,
		}),
	}
}

func (c *Consumer) ConsumerMessages(handle func(key, value []byte) error) {
	for {
		m, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("error reading message: %v", err)
			break
		}

		if err := handle(m.Key, m.Value); err != nil {
			log.Printf("error handling message: %v", err)
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
