package game

import (
	"whatthecard/logger"
)

const maxCardsPerPlayer = 20

// Game represents a game
type Game struct {
	RoomID           string
	Players          map[int]*Player
	DrawPile         *Pile
	DiscardPile      *Pile
	lastPlayerID     int
	HostID           int
	LastDrawPlayerID int
	CardsPerPlayer   int
	logger           *logger.Logger
}

// NewGame returns a new Game
func NewGame(logger *logger.Logger) *Game {
	return &Game{
		Players:        make(map[int]*Player),
		DrawPile:       NewPile(),
		DiscardPile:    NewPile(),
		CardsPerPlayer: 5,
		logger:         logger,
	}
}

// State represents a game state
type State struct {
	Phase            string    `json:"phase"`
	DrawPileLeft     int       `json:"draw_pile_left"`
	DiscardPile      *Pile     `json:"discard_pile"`
	Players          []*Player `json:"players"`
	LastDrawPlayerID int       `json:"last_draw_player_id"`
}

// SetCardsPerPlayer resets the game and sets a number of cards per player
func (g *Game) SetCardsPerPlayer(n int) {
	if n <= 0 || n > maxCardsPerPlayer {
		return
	}
	g.Reset(0)
	g.SetCardsPerPlayer(n)
}

// AddPlayer adds a player to the game
func (g *Game) AddPlayer(name string) *Player {
	id := g.lastPlayerID + 1
	p := NewPlayer(id, name)
	g.Players[p.ID] = p
	g.logger.Debugf("player %d has joined the room %s", id, g.RoomID)
	return p
}

// RemovePlayer removes a player from a game
func (g *Game) RemovePlayer(id int) {
	_, ok := g.Players[id]
	if !ok {
		return
	}
	delete(g.Players, id)
	g.logger.Debugf("player %d has left the room %s", id, g.RoomID)
}

// DrawCard draws a card
func (g *Game) DrawCard(playerID int) *Card {
	g.LastDrawPlayerID = playerID
	return g.DrawPile.Pop()
}

// Reset resets a game
// modes
// 0: delete all cards from piles
// 1: reset piles
func (g *Game) Reset(mode int) {
	switch mode {
	case 0:
		g.DrawPile.Reset()
		g.DiscardPile.Reset()
	case 1:
		g.DrawPile.Cards = append(g.DrawPile.Cards, g.DiscardPile.Cards...)
		g.DiscardPile.Reset()
	}
	for _, player := range g.Players {
		player.NumberOfSubmittedCards = 0
	}
}

// AddCard adds a card to the game
func (g *Game) AddCard(text string, playerID int) *Card {
	player, ok := g.Players[playerID]
	if !ok {
		return nil
	}

	card := &Card{
		Text:   text,
		Author: player.Name,
	}
	g.DrawPile.Push(card)
	player.NumberOfSubmittedCards++
	return card
}

// State returns a game state for player with given player id
func (g Game) State(playerID int) State {
	phase := "play"
	for _, player := range g.Players {
		if player.NumberOfSubmittedCards < g.CardsPerPlayer {
			phase = "submit_cards"
			break
		}
	}

	players := make([]*Player, 0, len(g.Players))
	for i := 1; i <= len(g.Players); i++ {
		player, ok := g.Players[i]
		if ok {
			players = append(players, player)
		}
	}
	return State{
		Phase:            phase,
		DrawPileLeft:     g.DrawPile.Len(),
		DiscardPile:      g.DiscardPile,
		Players:          players,
		LastDrawPlayerID: g.LastDrawPlayerID,
	}
}

// Command represents a game command
type Command struct {
	Name    string
	Payload interface{}
}

type (
	// SetCardPerPlayerPayload is a set cards per player payload
	SetCardPerPlayerPayload struct {
		CardsPerPlayer int
	}

	// RemovePlayerPayload is a remove player payload
	RemovePlayerPayload struct {
		ID int
	}

	// DrawCardPayload is a draw card payload
	DrawCardPayload struct {
		PlayerID int
	}

	// AddCardPayload is an add card payload
	AddCardPayload struct {
		Name     string
		PlayerID int
	}

	// ResetPayload is a reset payload
	ResetPayload struct {
		Mode int
	}
)

// ExecCommand executes a command
func (g *Game) ExecCommand(cmd Command) {
	switch cmd.Name {
	case "set_cards_per_player":
		payload, ok := cmd.Payload.(SetCardPerPlayerPayload)
		if !ok {
			return
		}
		g.SetCardsPerPlayer(payload.CardsPerPlayer)
	case "remove_player":
		payload, ok := cmd.Payload.(RemovePlayerPayload)
		if !ok {
			return
		}
		g.RemovePlayer(payload.ID)
	case "draw_card":
		payload, ok := cmd.Payload.(DrawCardPayload)
		if !ok {
			return
		}
		g.DrawCard(payload.PlayerID)
	case "add_card":
		payload, ok := cmd.Payload.(AddCardPayload)
		if !ok {
			return
		}
		g.AddCard(payload.Name, payload.PlayerID)
	case "reset":
		payload, ok := cmd.Payload.(ResetPayload)
		if !ok {
			return
		}
		g.Reset(payload.Mode)
	}
}
