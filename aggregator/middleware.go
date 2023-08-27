package main

import (
	"time"

	"github.com/gastrader/407ETR/types"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) Aggregator {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetOutput(colorable.NewColorableStdout())
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
		}).Info("AggregateDistance")
	}(time.Now())
	err = m.next.AggregateDistance(distance)
	return
}

func (m *LogMiddleware) CalculateInvoice(obuID int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		var (
			dist float64
			amt  float64
		)
		if inv != nil {
			dist = inv.TotalDistance
			amt = inv.TotalAmount
		}
		logrus.WithFields(logrus.Fields{
			"took":      time.Since(start),
			"err":       err,
			"obuID":     obuID,
			"totalDist": dist,
			"totalAmt":  amt,
		}).Info("Calc Invoice")
	}(time.Now())
	inv, err = m.next.CalculateInvoice(obuID)
	return
}
