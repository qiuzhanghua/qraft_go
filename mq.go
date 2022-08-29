package main

import (
	"github.com/segmentio/kafka-go"
	"os"
	"strings"
	"time"
)

var topic = "my-topic"
var mq *kafka.Writer

const DefaultKafkaAddress = "localhost:9092,localhost:9093,localhost:9094,localhost:9095,localhost:9096,localhost:9097"

func init() {
	kafkaAddress := os.Getenv("KAFKA_ADDRESS")
	if len(kafkaAddress) < 5 {
		kafkaAddress = DefaultKafkaAddress
	}

	addr := kafka.TCP(strings.Split(kafkaAddress, ",")...)
	mq = &kafka.Writer{
		Addr:         addr,
		Topic:        topic,
		BatchSize:    8192,
		BatchTimeout: time.Second * 2,
		Async:        true,
	}
}
