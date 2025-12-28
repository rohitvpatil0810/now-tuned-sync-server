package main

import (
	"log"
	"net/http"

	"github.com/rohitvpatil0810/now-tuned-sync-server/internal/handler"
	"github.com/rohitvpatil0810/now-tuned-sync-server/internal/store"
	"github.com/rs/cors"
)

func main() {
	mux := http.NewServeMux()

	store := store.NewPartitionedMusicState()
	musicStateHandler := handler.NewMusicStateHandler(store)

	c := cors.AllowAll()

	addr := ":8080"
	log.Println("now-tuned-sync-server listening on", addr)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello From now-tuned-sync-server"))
	})

	mux.HandleFunc("/music-state", musicStateHandler.UpdateMusicState)
	mux.HandleFunc("/music-state/winner", musicStateHandler.GetWinnerMusicState)
	mux.HandleFunc("/music-state/delete", musicStateHandler.DeleteMusicState)

	if err := http.ListenAndServe(addr, c.Handler(mux)); err != nil {
		log.Fatal(err)
	}
}
