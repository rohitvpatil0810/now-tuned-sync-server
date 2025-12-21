package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	addr := ":8080"
	log.Println("now-tuned-sync-server listening on", addr)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello From now-tuned-sync-server"))
	})

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
