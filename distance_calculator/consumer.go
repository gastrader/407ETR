package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gastrader/407ETR/aggregator/client"
	"github.com/gastrader/407ETR/types"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

//OR KAFKA TRANSPORT
type KafkaConsumer struct {
	consumer  *kafka.Consumer
	isRunning bool
	calcService CalculatorServicer
	aggClient client.Client
}

func NewKafkaConsumer(topic string, svc CalculatorServicer, aggClient client.Client) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",

	})
	if err != nil {
		return nil, err
	}

	c.SubscribeTopics([]string{topic}, nil)
	return &KafkaConsumer{
		consumer: c,
		calcService: svc,
		aggClient: aggClient,
	}, nil
}

func (c *KafkaConsumer) Start() {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetOutput(colorable.NewColorableStdout())
	logrus.Info("kafka transport started")
	c.isRunning = true
	c.readMessageLoop()
}

func (c *KafkaConsumer) readMessageLoop() {
	for c.isRunning {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("kafka consume error %s", err)
			continue
		}
		var data types.OBUData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("JSON serialization error %s", err)
			continue
		}
		distance, err := c.calcService.CalculateDistance(data)
		if err != nil {
			logrus.Errorf("Calculation error %s", err)
			continue
		}
		req := &types.AggregateRequest{
			Value: distance,
			Unix: time.Now().UnixNano(),
			ObuID: int32(data.OBUID),
		}
		if err := c.aggClient.Aggregate(context.Background(), req); err != nil {
			log.Fatal(err)
			continue
		}
	}
}

