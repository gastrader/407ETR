package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gastrader/407ETR/aggregator/client"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func main() {
	listenAddr := flag.String("listenAddr", ":6000", "Listen address for GATEWAY")
	aggregatorServiceAddr := flag.String("aggServiceAddr", "127.0.0.1:3000", "Listen address for aggregator service")
	flag.Parse()
	var (
		client     = client.NewHTTPClient(*aggregatorServiceAddr) //endpoint of aggregator service
		invHandler = newInvoiceHandler(client)
	)
	http.HandleFunc("/invoice", makeAPIFunc(invHandler.handleGetInvoice))
	fmt.Printf("gateway HTTP server running on port %s\n", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))

}

type InvoiceHandler struct {
	client client.Client
}

func newInvoiceHandler(c client.Client) *InvoiceHandler {
	return &InvoiceHandler{client: c}
}

func (h *InvoiceHandler) handleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	inv, err := h.client.GetInvoice(context.Background(), 1554080224)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, inv)
}

func writeJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func makeAPIFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
		logrus.SetOutput(colorable.NewColorableStdout())
		defer func(start time.Time) {
			logrus.WithFields(logrus.Fields{
				"took": time.Since(start),
				"uri":  r.RequestURI,
			}).Info("REQ ::  ")
		}(time.Now())
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
}
