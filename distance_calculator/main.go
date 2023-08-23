package main

import "log"

type DistanceCalculator struct {
}

const kafkaTopic = "obudata"

//Transport could be HHTP, GRPC, KAFKA -> attach business logic to transport.

func main() {
	var (
		err error
		svc CalculatorServicer
	)
	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}