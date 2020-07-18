package game

import (
	"math/rand"
	"time"
)

// Card represents a card
type Card struct {
	ID     int    `json:"id"`
	Text   string `json:"text"`
	Author string `json:"author"`
}

// NewCard returns a new Card
func NewCard(id int, text, author string) *Card {
	return &Card{
		ID:     id,
		Text:   text,
		Author: author,
	}
}

// Pile represents a card pile
type Pile struct {
	Cards      []*Card `json:"cards"`
	lastCardID int
}

// NewPile returns an empty Pile
func NewPile() *Pile {
	p := &Pile{}
	p.Reset()
	return p
}

// Reset removes all cards from the pile
func (p *Pile) Reset() {
	p.Cards = make([]*Card, 0)
	p.lastCardID = 0
}

// Shuffle shuffles cards in the pile
func (p *Pile) Shuffle() {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	r.Shuffle(len(p.Cards), func(i, j int) { p.Cards[i], p.Cards[j] = p.Cards[j], p.Cards[i] })
}

// Len return a length of the pile
func (p Pile) Len() int {
	return len(p.Cards)
}

// Pop pops the top card from the pile
func (p *Pile) Pop() *Card {
	if p.Len() == 0 {
		return nil
	}
	card := p.Cards[p.Len()-1]
	p.Cards = p.Cards[:p.Len()-1]
	return card
}

// Push pushs a card to the top of the pile
func (p *Pile) Push(card *Card) {
	p.lastCardID++
	card.ID = p.lastCardID
	p.Cards = append(p.Cards, card)
}
