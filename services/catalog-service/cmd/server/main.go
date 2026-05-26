package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("CATALOG_SERVICE_PORT")
	if port == "" {
		port = "8082"
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("catalog-service: ok"))
	})

	addr := ":" + port
	log.Printf("catalog-service listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
