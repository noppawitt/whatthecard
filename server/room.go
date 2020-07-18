package server

import (
	"encoding/json"
	"fmt"
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
	r.lastClientID++
	id := r.lastClientID
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
				r.logger.Error(err)
				return
			}
		case msgByte := <-client.send:
			if len(msgByte) == 0 {
				break
			}

			msg := &Message{}
			err := json.Unmarshal(msgByte, msg)
			if err != nil {
				r.logger.Error(err)
				break
			}
			cmd, err := msg.ToGameCommand(clientID)
			if err != nil {
				r.logger.Error(err)
				break
			}

			if err = r.game.ExecCommand(cmd); err != nil {
				r.logger.Error(err)
				break
			}

			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			r.BroadcastState()
		}
	}
}

// BroadcastState broadcasts latest game state to all clients
func (r *Room) BroadcastState() {
	for _, player := range r.game.Players {
		state := r.game.State(player.ID)
		client := r.clients[player.ID]
		if client != nil {
			client.WriteJSON(state)
		}
	}
}

// Message represents a message from the websocket client
type Message struct {
	Name    string          `json:"name"`
	Payload json.RawMessage `json:"payload"`
}

// ToGameCommand converts a Message to a Game Command
func (m Message) ToGameCommand(playerID int) (game.Command, error) {
	cmd := game.Command{
		Name:     m.Name,
		PlayerID: playerID,
	}
	var payload interface{}
	switch cmd.Name {
	case "set_cards_per_player":
		payload = &game.SetCardPerPlayerPayload{}
	case "add_player":
		payload = &game.AddPlayerPayload{}
	case "remove_player":
		payload = &game.RemovePlayerPayload{}
	case "add_card":
		payload = &game.AddCardPayload{}
	case "reset":
		payload = &game.ResetPayload{}
	case "start", "draw_card":
		return cmd, nil
	default:
		return cmd, fmt.Errorf("invalid game command: %s", cmd.Name)
	}
	err := json.Unmarshal(m.Payload, payload)
	if err != nil {
		return cmd, err
	}
	cmd.Payload = payload
	return cmd, nil
}
