package game

// Player represents a player
type Player struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	HasSubmitted bool   `json:"has_submitted"`
}

// NewPlayer returns a new Player
func NewPlayer(id int, name string) *Player {
	return &Player{
		ID:           id,
		Name:         name,
		HasSubmitted: false,
	}
}
