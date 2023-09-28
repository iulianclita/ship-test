package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iulianclita/ship-test/ship"
)

const (
	orderQtysQueyParam  = "order_qty"
	packSizesQueryParam = "pack_sizes"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ship", shipPacks)

	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to listen: %s\n", err)
		}
	}()
	log.Print("Server Started")

	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shutdown the http server: %v", err)
	}
	log.Print("Server Shutdown Successfully")
}

type response struct {
	Data  map[int]int `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func shipPacks(w http.ResponseWriter, r *http.Request) {
	orderQty, err := ship.ExtractOrderQty(r.URL.Query().Get(orderQtysQueyParam))
	if err != nil {
		switch {
		case errors.Is(err, ship.ErrOrderQtyInvalidFormat):
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(response{
				Error: ship.ErrOrderQtyInvalidFormat.Error(),
			}); err != nil {
				log.Printf("failed to encode http response: %v", err)
			}
		case errors.Is(err, ship.ErrOrderQtyInvalidValue):
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(response{
				Error: ship.ErrOrderQtyInvalidValue.Error(),
			}); err != nil {
				log.Printf("failed to encode http response: %v", err)
			}
		default:
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(response{
				Error: err.Error(),
			}); err != nil {
				log.Printf("failed to encode http response: %v", err)
			}
		}
		return
	}

	packSizes, err := ship.ExtractPackSizes(r.URL.Query().Get(packSizesQueryParam))
	if err != nil {
		switch {
		case errors.Is(err, ship.ErrPackSizesInvalidFormat):
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(response{
				Error: ship.ErrPackSizesInvalidFormat.Error(),
			}); err != nil {
				log.Printf("failed to encode http response: %v", err)
			}
		case errors.Is(err, ship.ErrPackSizeInvalidValue):
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(response{
				Error: ship.ErrPackSizeInvalidValue.Error(),
			}); err != nil {
				log.Printf("failed to encode http response: %v", err)
			}
		default:
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(response{
				Error: err.Error(),
			}); err != nil {
				log.Printf("failed to encode http response: %v", err)
			}
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response{
		Data: ship.CalculatePacksToShip(orderQty, packSizes),
	}); err != nil {
		log.Printf("failed to encode http response: %v", err)
	}

}
