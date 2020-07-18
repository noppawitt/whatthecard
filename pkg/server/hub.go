package server

import (
	"math/rand"
	"net/http"
	"time"
	"whatthecard/pkg/game"
	"whatthecard/pkg/logger"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// RoomIDLength is the length of a room id
const RoomIDLength = 4

// Hub is central handler for websocket connections
type Hub struct {
	rooms    map[string]*Room
	upgrader websocket.Upgrader
	logger   *logger.Logger
}

// NewHub returns a new Hub
func NewHub(logger *logger.Logger) *Hub {
	return &Hub{
		rooms: make(map[string]*Room),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(*http.Request) bool { return true },
		},
		logger: logger,
	}
}

// GetRoom returns a room with the given room id
func (h Hub) GetRoom(id string) *Room {
	return h.rooms[id]
}

// CreateRoom creates a new room
func (h *Hub) CreateRoom(game *game.Game, logger *logger.Logger) *Room {
	for {
		id := randString(RoomIDLength)
		_, ok := h.rooms[id]
		if !ok {
			h.rooms[id] = NewRoom(id, game, logger)
			game.RoomID = id
			h.logger.Debugf("room %s has been created", id)
			return h.rooms[id]
		}
	}
}

func randString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return string(b)
}

// HandleWS handles websocket connection
func (h *Hub) HandleWS(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID := vars["id"]
	if roomID == "" {
		writeError(w, "room id is required", http.StatusBadRequest)
	}

	playerName := r.URL.Query().Get("player_name")
	if playerName == "" {
		writeError(w, "player_name is required", http.StatusBadRequest)
		return
	}

	room := h.GetRoom(roomID)
	if room == nil {
		h.logger.Debug("room not found:", roomID)
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error(err)
		return
	}

	client := NewClient(conn, h.logger)
	clientID := room.Join(client)
	player := room.game.AddPlayer(clientID, playerName)
	room.BroadcastState()

	go room.WritePump(clientID)
	client.ReadPump()

	room.Leave(clientID)
	room.game.RemovePlayer(player.ID)
	room.BroadcastState()

	if room.TotalClient == 0 {
		delete(h.rooms, room.ID)
		h.logger.Debugf("room %s has been deleted", room.ID)
	}
}
