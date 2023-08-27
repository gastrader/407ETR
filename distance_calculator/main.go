package main

import (
	"log"

	"github.com/gastrader/407ETR/aggregator/client"
)

type DistanceCalculator struct {
}

const (
	kafkaTopic = "obudata"
	aggregatorEndpoint = "http://127.0.0.1:3000/aggregate"
)

//Transport could be HHTP, GRPC, KAFKA -> attach business logic to transport.

func main() {
	var (
		err error
		svc CalculatorServicer
	)
	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc, client.NewClient(aggregatorEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}