package game

import "whatthecard/pkg/logger"

// Service represents a game service
type Service struct {
	logger *logger.Logger
}

// NewService returns a new GameService
func NewService(logger *logger.Logger) *Service {
	return &Service{
		logger: logger,
	}
}

// NewGame returns a new Game
func (s *Service) NewGame() *Game {
	return NewGame(s.logger)
}
