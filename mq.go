package main

import (
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

var topic = "my-topic"
var mq *kafka.Writer

const DefaultKafkaAddress = "localhost:9092,localhost:9093,localhost:9094,localhost:9095,localhost:9096,localhost:9097"

func KafkaInit() {
	viper.SetDefault("kafka.cluster", DefaultKafkaAddress)
	kafkaAddress := viper.GetString("kafka.cluster")
	addresses := os.Getenv("KAFKA_CLUSTER")
	if len(addresses) >= 5 {
		kafkaAddress = addresses
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
