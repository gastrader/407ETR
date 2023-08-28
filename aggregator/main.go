package main

import (
	// "context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	// "time"

	// "github.com/gastrader/407ETR/aggregator/client"
	"github.com/gastrader/407ETR/types"
	"google.golang.org/grpc"
)

func main() {
	httpLlistenAddr := flag.String("httpAddr", ":3000", "listen address of HTTP server")
	grpcLlistenAddr := flag.String("grpcAddr", "localhost:3001", "listen address of GRPC server")

	flag.Parse()

	store := NewMemoryStore()
	var (
		svc Aggregator
	)
	svc = NewInvoiceAggregator(store)
	svc = NewLogMiddleware(svc)
	go func() {
		log.Fatal(makeGRPCTransport(*grpcLlistenAddr, svc))
	}()
	log.Fatal(makeHTTPTransport(*httpLlistenAddr, svc))

}

func makeHTTPTransport(listenAddr string, svc Aggregator) error {
	fmt.Println("HTTP transport running on port ", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleGetInvoice(svc))
	return http.ListenAndServe(listenAddr, nil)
}

func makeGRPCTransport(listenAddr string, svc Aggregator) error {
	fmt.Println("gRPC trasnport running on port: ", listenAddr)
	// Make TCP Listener
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	// Make GRPC native server with options
	server := grpc.NewServer([]grpc.ServerOption{}...)
	// Register GRPC server implementation to the GRPC package
	types.RegisterAggregatorServer(server, NewAggregatorGRPCServer(svc))
	return server.Serve(ln)
}

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query().Get("obu")
		if values == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing OBU ID"})
			return
		}
		obuID, err := strconv.Atoi(values)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad ID"})
			return
		}
		invoice, err := svc.CalculateInvoice(obuID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, invoice)

	}
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := svc.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
