package main

import (
	"time"

	"github.com/gastrader/407ETR/types"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next CalculatorServicer
}

func NewLogMiddleware(next CalculatorServicer) CalculatorServicer {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetOutput(colorable.NewColorableStdout())
	return &LogMiddleware{
		next: next,
	}
}
func (m *LogMiddleware) CalculateDistance(data types.OBUData) (dist float64, err error){
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err": err,
			"dist": dist,
			"id": data.OBUID,
		}).Info("calculate distance")
	}(time.Now())
	dist, err = m.next.CalculateDistance(data)
	return dist, err
}