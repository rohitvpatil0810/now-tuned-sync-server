package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/rohitvpatil0810/now-tuned-sync-server/internal/model"
	"github.com/rohitvpatil0810/now-tuned-sync-server/internal/store"
)

type MusicStateHandler struct {
	store *store.PartitionedMusicState
}

func NewMusicStateHandler(store *store.PartitionedMusicState) *MusicStateHandler {
	return &MusicStateHandler{
		store: store,
	}
}

func (h *MusicStateHandler) UpdateMusicState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.MusicState

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Println("Incoming request state 2 : ", req)

	h.store.Update(req)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *MusicStateHandler) GetWinnerMusicState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	winner := h.store.Winner()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(winner)
}
