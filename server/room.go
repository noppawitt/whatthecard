package server

import (
	"encoding/json"
	"time"
	"whatthecard/game"
	"whatthecard/logger"

	"github.com/gorilla/websocket"
)

// Room represents a client room
type Room struct {
	ID           string
	clients      map[int]*Client
	lastClientID int
	TotalClient  int
	game         *game.Game
	logger       *logger.Logger
}

// NewRoom returns a new Room
func NewRoom(id string, game *game.Game, logger *logger.Logger) *Room {
	return &Room{
		ID:           id,
		clients:      make(map[int]*Client, 0),
		lastClientID: 0,
		TotalClient:  0,
		game:         game,
		logger:       logger,
	}
}

// Join joins the client to the room
func (r *Room) Join(client *Client) int {
	id := r.lastClientID + 1
	r.clients[id] = client
	r.TotalClient++
	return id
}

// Leave leaves the client from the room
func (r *Room) Leave(clientID int) {
	delete(r.clients, clientID)
	r.TotalClient--
}

// WritePump writes messages to the connected clients and checks if any clients are disconnected
func (r *Room) WritePump(clientID int) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		client := r.clients[clientID]
		if client == nil {
			break
		}

		select {
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case msgByte := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			msg := &Message{}
			err := json.Unmarshal(msgByte, msg)
			if err != nil {
				return
			}
			cmd, err := msg.ToGameCommand()
			if err != nil {
				return
			}

			r.game.ExecCommand(cmd)

			for _, player := range r.game.Players {
				state := r.game.State(player.ID)
				client := r.clients[player.ID]
				if client != nil {
					client.WriteJSON(state)
				}
			}
		}
	}
}

// Message represents a message from the websocket client
type Message struct {
	Name    string          `json:"name"`
	Payload json.RawMessage `json:"payload"`
}

// ToGameCommand converts a Message to a Game Command
func (m Message) ToGameCommand() (game.Command, error) {
	cmd := game.Command{}
	var payload interface{}
	switch m.Name {
	case "remove_player":
		payload = game.RemovePlayerPayload{}
	case "draw_card":
		payload = game.DrawCardPayload{}
	case "add_card":
		payload = game.AddCardPayload{}
	case "reset":
		payload = game.ResetPayload{}
	}
	err := json.Unmarshal(m.Payload, &payload)
	if err != nil {
		return cmd, err
	}
	cmd.Payload = payload
	return cmd, nil
}
