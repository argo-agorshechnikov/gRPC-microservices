package kafka

import "os"

func GetBrokers() []string {
	broker := os.Getenv("KAFKA_BROKER_ADDRESS")

	if broker == "" {
		broker = "kafka:9092"
	}

	return []string{broker}
}
