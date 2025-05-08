package main

import (
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/whywaita/poker-go"
)

var _ tea.Model = &Cards{}

const spaceBetweenCardsWithContent = 4

type Cards struct {
	allCards     []Card
	hasChanged   bool
	selectedCard *Card
}

func (c *Cards) HasChanged() bool {
	if c.hasChanged {
		c.hasChanged = false
		return true
	}
	return false
}

func NewCards() *Cards {
	m := Cards{}
	m.allCards = make([]Card, 52)
	newDeck := poker.NewDeck()
	for i := range len(newDeck.Cards) {
		c := NewCard(newDeck.DrawCard())

		m.allCards[i] = c
	}

	return &m
}

// Init implements tea.Model.
func (c Cards) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (c *Cards) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft {
			log.Println("Mouse clicked: %v, %v", msg.Action, msg.Button)
			ok, card := c.GetCardByCoordinates(msg)
			if ok {
				c.hasChanged = true
				c.selectedCard = card
			}
		}
	}

	return c, nil
}

func (c Cards) GetCardByCoordinates(mouseMsg tea.MouseMsg) (bool, *Card) {
	if mouseMsg.X < 0 {
		return false, nil
	}
	if mouseMsg.X >= 52 {
		return false, nil
	}

	if mouseMsg.Y < 0 {
		return false, nil
	}
	if mouseMsg.Y >= 4 {
		return false, nil
	}

	moduloX := mouseMsg.X % spaceBetweenCardsWithContent
	if moduloX == 2 || moduloX == 3 {
		return false, nil
	}

	// x and y are matricial view
	x := mouseMsg.X / spaceBetweenCardsWithContent
	y := mouseMsg.Y

	card := c.allCards[y*13+x]

	log.Printf("Card selected: %s", card.String())
	return true, &card
}

// View implements tea.Model.
func (c Cards) View() string {
	sb := strings.Builder{}
	for row := range 4 {
		for col := range 13 {
			card := c.allCards[row*13+col]
			sb.WriteString((card.String()))
			sb.WriteString("  ")
		}
		if row != 3 {
		sb.WriteString("\n")
		}
	}
	return sb.String()
}
