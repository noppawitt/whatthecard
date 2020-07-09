package server

import (
	"log"
	"math/rand"
	"net/http"
	"whatthecard/game"
	"whatthecard/logger"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const RoomIDLength = 4

type Hub struct {
	rooms    map[string]*Room
	upgrader websocket.Upgrader
	logger   *logger.Logger
}

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

func (h Hub) GetRoom(id string) *Room {
	return h.rooms[id]
}

func (h *Hub) CreateRoom(game *game.Game, logger *logger.Logger) *Room {
	for {
		id := randString(RoomIDLength)
		_, ok := h.rooms[id]
		if !ok {
			h.rooms[id] = NewRoom(id, game, logger)
			game.RoomID = id
			h.logger.Debugf("room %s has beed created", id)
			return h.rooms[id]
		}
	}
}

func randString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (h *Hub) HandleWs(w http.ResponseWriter, r *http.Request) {
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
		log.Println("room not found:", roomID)
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(conn)
	clientID := room.Join(client)
	player := room.game.AddPlayer(playerName)

	go room.WritePump(clientID)
	client.ReadPump()

	room.Leave(clientID)
	room.game.RemovePlayer(player.ID)
	if room.TotalClient == 0 {
		delete(h.rooms, room.ID)
		h.logger.Debugf("room %s has beed deleted", room.ID)
	}
}
