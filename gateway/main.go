package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gastrader/407ETR/aggregator/client"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func main() {
	listenAddr := flag.String("listenAddr", ":6000", "Listen address for GATEWAY")
	flag.Parse()
	var (
		client     = client.NewHTTPClient("TODO")
		invHandler = newInvoiceHandler(client)
	)
	http.HandleFunc("/invoice", makeAPIFunc(invHandler.HandleGetInvoice))
	fmt.Printf("gateway HTTP server running on port %s", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))

}

type InvoiceHandler struct {
	client client.Client
}

func newInvoiceHandler(c client.Client) *InvoiceHandler {
	return &InvoiceHandler{client: c}
}

func (h *InvoiceHandler) HandleGetInvoice(w http.ResponseWriter, r *http.Request) error {

	inv, err := h.client.GetInvoice(context.Background(), 696969)
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
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
}
