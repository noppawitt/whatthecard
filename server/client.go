package server

import (
	"time"
	"whatthecard/logger"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

// Client provides a websocket connection
type Client struct {
	conn   *websocket.Conn
	send   chan []byte
	logger *logger.Logger
}

// NewClient returns a new Client
func NewClient(conn *websocket.Conn, logger *logger.Logger) *Client {
	return &Client{
		conn:   conn,
		logger: logger,
		send:   make(chan []byte),
	}
}

// ReadPump reads for an incomming message
func (c *Client) ReadPump() {
	defer close(c.send)

	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		c.logger.Debug(string(message))
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			c.logger.Error(err)
			break
		}
		c.send <- message
	}
}

// WriteJSON writes a JSON to the client
func (c *Client) WriteJSON(v interface{}) {
	err := c.conn.WriteJSON(v)
	if err != nil {
		c.logger.Error(err)
	}
}
