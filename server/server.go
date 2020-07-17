package server

import (
	"encoding/json"
	"net/http"
	"whatthecard/game"
	"whatthecard/logger"

	"github.com/gorilla/mux"
)

// Server represents a server
type Server struct {
	r           *mux.Router
	hub         *Hub
	gameService *game.Service
	logger      *logger.Logger
}

// New returns a new Server
func New(hub *Hub, gameService *game.Service, logger *logger.Logger) *Server {
	return &Server{
		r:           mux.NewRouter(),
		hub:         hub,
		gameService: gameService,
		logger:      logger,
	}
}

// Start starts the server
func (s *Server) Start(addr string) error {
	s.registerRoutes()
	http.Handle("/", s.r)
	s.logger.Debugf("server is running at %s", addr)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) registerRoutes() {
	s.r.HandleFunc("/", s.handleHome).Methods(http.MethodGet)
	s.r.HandleFunc("/room", s.handleCreateRoom).Methods(http.MethodPost)
	s.r.HandleFunc("/ws/room/{id}", s.hub.HandleWS).Methods(http.MethodGet)
}

func (Server) handleHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

func (s *Server) handleCreateRoom(w http.ResponseWriter, r *http.Request) {
	game := s.gameService.NewGame()
	room := s.hub.CreateRoom(game, s.logger)

	writeJSON(w, map[string]interface{}{"room_id": room.ID})
}

func writeJSON(w http.ResponseWriter, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, err string, code int) error {
	w.WriteHeader(code)
	return writeJSON(w, map[string]interface{}{
		"error": err,
	})
}
