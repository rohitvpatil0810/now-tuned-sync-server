package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

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

	// update timestamp using time in req
	req.UpdatedAt = time.Now()

	log.Printf("Updating music state for clientId: %s, playbackState: %s\n", req.Client.ClientID, req.PlaybackState)

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

	if winner == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(nil)
		return
	}

	// Create response without client field
	response := map[string]interface{}{
		"metadata":      winner.Metadata,
		"playbackState": winner.PlaybackState,
		"updatedAt":     winner.UpdatedAt,
		"version":       winner.Version,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *MusicStateHandler) DeleteMusicState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	clientID := r.URL.Query().Get("clientId")
	if clientID == "" {
		http.Error(w, "clientId query parameter is required", http.StatusBadRequest)
		return
	}

	log.Println("Deleting music state for clientId:", clientID)
	h.store.Remove(clientID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}
