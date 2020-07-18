package server

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"whatthecard/pkg/game"
	"whatthecard/pkg/logger"

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
	s.r.HandleFunc("/room", s.handleCreateRoom).Methods(http.MethodPost)
	s.r.HandleFunc("/ws/room/{id}", s.hub.HandleWS).Methods(http.MethodGet)

	spa := spaHandler{staticPath: "./web/dist", indexPath: "index.html"}
	s.r.PathPrefix("/").Handler(spa)
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

// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type spaHandler struct {
	staticPath string
	indexPath  string
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}
