package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rohitvpatil0810/now-tuned-sync-server/internal/model"
	"github.com/rohitvpatil0810/now-tuned-sync-server/internal/store"
)

type MusicStateHandler struct {
	store *store.PartitionedMusicState
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins - for testing purposes
	},
}

type WSClient struct {
	conn *websocket.Conn
	send chan *model.MusicState
}

func NewMusicStateHandler(store *store.PartitionedMusicState) *MusicStateHandler {
	return &MusicStateHandler{
		store: store,
	}
}

// toResponseFormat converts MusicState to response format without client data
func toResponseFormat(state *model.MusicState) map[string]interface{} {
	if state == nil {
		return nil
	}
	return map[string]interface{}{
		"metadata":      state.Metadata,
		"playbackState": state.PlaybackState,
		"updatedAt":     state.UpdatedAt,
		"version":       state.Version,
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
	response := toResponseFormat(winner)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *MusicStateHandler) WinnerWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Println("New WebSocket connection")

	// Send current winner immediately upon connection
	winner := h.store.Winner()
	response := toResponseFormat(winner)
	if response != nil {
		if err := conn.WriteJSON(response); err != nil {
			log.Println("WebSocket write error:", err)
			return
		}
	}

	// Register client to receive updates
	client := &WSClient{
		conn: conn,
		send: make(chan *model.MusicState, 10),
	}
	h.store.RegisterClient(client.send)
	defer h.store.UnregisterClient(client.send)

	// Listen for updates and send to client
	for state := range client.send {
		response := toResponseFormat(state)
		if err := client.conn.WriteJSON(response); err != nil {
			log.Println("WebSocket write error:", err)
			return
		}
	}
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
