package game

import (
	"fmt"
	"math"
	"whatthecard/pkg/logger"
)

const maxCardsPerPlayer = 20

// Phase represents a game phase
type Phase int

// Game phases
const (
	WaitingPhase Phase = iota
	SubmitPhase
	PlayPhase
)

func (p Phase) String() string {
	switch p {
	case WaitingPhase:
		return "WAITING_PHASE"
	case SubmitPhase:
		return "SUBMIT_PHASE"
	case PlayPhase:
		return "PLAY_PHASE"
	default:
		return ""
	}
}

// Game represents a game
type Game struct {
	RoomID           string
	Phase            Phase
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
		Phase:          WaitingPhase,
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
	DiscardCards     []*Card   `json:"discard_cards"`
	CardsPerPlayer   int       `json:"cards_per_player"`
	PlayerID         int       `json:"player_id"`
	HostID           int       `json:"host_id"`
	Players          []*Player `json:"players"`
	LastDrawPlayerID int       `json:"last_draw_player_id"`
}

// SetCardsPerPlayer resets the game and sets a number of cards per player
func (g *Game) SetCardsPerPlayer(n int) {
	if n <= 0 || n > maxCardsPerPlayer {
		return
	}
	g.CardsPerPlayer = n
}

// AddPlayer adds a player to the game
func (g *Game) AddPlayer(id int, name string) *Player {
	p := NewPlayer(id, name)
	g.Players[p.ID] = p
	g.logger.Debugf("player %d has joined the room %s", id, g.RoomID)

	if len(g.Players) == 1 {
		g.PromoteHost(p.ID)
	}

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

	if len(g.Players) > 0 && g.isHost(id) {
		minID := math.MaxInt32
		for k := range g.Players {
			if k < minID && g.Players[k] != nil {
				minID = k
			}
		}
		g.PromoteHost(minID)
	}
}

// PromoteHost promotes a player to a host
func (g *Game) PromoteHost(playerID int) {
	g.HostID = playerID
	g.logger.Debugf("player %d has been promoted to a host", playerID)
}

func (g *Game) isHost(playerID int) bool {
	return playerID == g.HostID
}

// DrawCard draws a card
func (g *Game) DrawCard(playerID int) *Card {
	g.LastDrawPlayerID = playerID
	card := g.DrawPile.Pop()
	g.DiscardPile.Cards = append(g.DiscardPile.Cards, card)
	return card
}

// Reset resets a game
// modes
// 0: delete all cards from piles then go back to waiting phase
// 1: reset piles
func (g *Game) Reset(mode int) {
	switch mode {
	case 0:
		g.DrawPile.Reset()
		g.DiscardPile.Reset()
		g.Phase = WaitingPhase
		for _, player := range g.Players {
			player.NumberOfSubmittedCards = 0
		}
	case 1:
		g.DrawPile.Cards = append(g.DrawPile.Cards, g.DiscardPile.Cards...)
		g.DiscardPile.Reset()
		g.DrawPile.Shuffle()
	}
}

// Start changes game phase to submit phase
func (g *Game) Start() {
	g.Phase = SubmitPhase
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

	allSubmitted := true
	for _, player := range g.Players {
		if player.NumberOfSubmittedCards < g.CardsPerPlayer {
			allSubmitted = false
			break
		}
	}

	if allSubmitted {
		g.DrawPile.Shuffle()
		g.Phase = PlayPhase
	}

	return card
}

// State returns a game state for player with given player id
func (g Game) State(playerID int) State {
	players := make([]*Player, 0, len(g.Players))
	for i := 1; i <= len(g.Players); i++ {
		player, ok := g.Players[i]
		if ok {
			players = append(players, player)
		}
	}
	return State{
		Phase:            g.Phase.String(),
		DrawPileLeft:     g.DrawPile.Len(),
		DiscardCards:     g.DiscardPile.Cards,
		CardsPerPlayer:   g.CardsPerPlayer,
		PlayerID:         playerID,
		HostID:           g.HostID,
		Players:          players,
		LastDrawPlayerID: g.LastDrawPlayerID,
	}
}

// Command represents a game command
type Command struct {
	Name     string
	PlayerID int
	Payload  interface{}
}

type (
	// SetCardPerPlayerPayload is a set cards per player payload
	SetCardPerPlayerPayload struct {
		CardsPerPlayer int `json:"cards_per_player"`
	}

	// AddPlayerPayload is an add player payload
	AddPlayerPayload struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	// RemovePlayerPayload is a remove player payload
	RemovePlayerPayload struct {
		ID int `json:"id"`
	}

	// AddCardPayload is an add card payload
	AddCardPayload struct {
		Text string `json:"text"`
	}

	// ResetPayload is a reset payload
	ResetPayload struct {
		Mode int `json:"mode"`
	}
)

// ExecCommand executes a command
func (g *Game) ExecCommand(cmd Command) error {
	switch cmd.Name {
	case "set_cards_per_player":
		if !g.isHost(cmd.PlayerID) {
			return CommandIsForHostOnlyErr{cmd: cmd}
		}
		payload, ok := cmd.Payload.(*SetCardPerPlayerPayload)
		if !ok {
			return InvalidCommandErr{cmd: cmd}
		}
		g.SetCardsPerPlayer(payload.CardsPerPlayer)
	case "add_player":
		payload, ok := cmd.Payload.(*AddPlayerPayload)
		if !ok {
			return InvalidCommandErr{cmd: cmd}
		}
		g.AddPlayer(payload.ID, payload.Name)
	case "remove_player":
		payload, ok := cmd.Payload.(*RemovePlayerPayload)
		if !ok {
			return InvalidCommandErr{cmd: cmd}
		}
		g.RemovePlayer(payload.ID)
	case "start":
		if !g.isHost(cmd.PlayerID) {
			return CommandIsForHostOnlyErr{cmd: cmd}
		}
		g.Start()
	case "draw_card":
		g.DrawCard(cmd.PlayerID)
	case "add_card":
		payload, ok := cmd.Payload.(*AddCardPayload)
		if !ok {
			return InvalidCommandErr{cmd: cmd}
		}
		g.AddCard(payload.Text, cmd.PlayerID)
	case "reset":
		if !g.isHost(cmd.PlayerID) {
			return CommandIsForHostOnlyErr{cmd: cmd}
		}
		payload, ok := cmd.Payload.(*ResetPayload)
		if !ok {
			return InvalidCommandErr{cmd: cmd}
		}
		g.Reset(payload.Mode)
	}

	return nil
}

// InvalidCommandErr occurs when a game command is invalid
type InvalidCommandErr struct {
	cmd Command
}

func (e InvalidCommandErr) Error() string {
	return fmt.Sprintf("cmd: %s with payload: %v is invalid", e.cmd.Name, e.cmd.Payload)
}

// CommandIsForHostOnlyErr occurs when a non host player try to execute host only command
type CommandIsForHostOnlyErr struct {
	cmd Command
}

func (e CommandIsForHostOnlyErr) Error() string {
	return fmt.Sprintf("player %d is not a host, cmd: %s must be executed by host player", e.cmd.PlayerID, e.cmd.Name)
}
