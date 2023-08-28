package main

import (
	"log"

	"github.com/gastrader/407ETR/aggregator/client"
)

type DistanceCalculator struct {
}

const (
	kafkaTopic = "obudata"
	aggregatorEndpoint = "127.0.0.1:3000/aggregate"
)



func main() {
	var (
		err error
		svc CalculatorServicer
	)
	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)

	// httpClient := client.NewHTTPClient(aggregatorEndpoint)
	grpcClient, err := client.NewGRPCClient(aggregatorEndpoint)
	if err != nil{
		log.Fatal(err)
	}

	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc, grpcClient)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}