package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"github.com/whywaita/poker-go"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type Model struct {
	list list.Model
	// allCards []Card
	players []poker.Player

	h1 [2]poker.Card
	h2 [2]poker.Card

	flop  [3]poker.Card
	turn  poker.Card
	river poker.Card

	nextSelection selection
}

type selection int

const (
	sp11 selection = iota
	sp12
	sp21
	sp22
	sf1
	sf2
	sf3
	st
	sr
)

func getNextSelection(s selection) selection {
	if s == 8 {
		return sp11
	}

	return s + 1
}

func NewModel() Model {
	m := Model{}
	// m.allCards = make([]Card, 52)
	newDeck := poker.NewDeck()
	li := make([]list.Item, 52)
	for i := range 52 {
		c := Card(newDeck.DrawCard())
		// m.allCards[i] = c
		li[i] = c
	}

	m.list = list.New(li, list.NewDefaultDelegate(), 0, 0)
	m.list.Title = "Poker Hand"

	return m
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View implements tea.Model.
func (m Model) View() string {
	return zone.Scan(docStyle.Render(m.list.View()))
}

type Card poker.Card

func (c Card) String() string {
	pc := poker.Card(c)
	return pc.Rank.String() + pc.Suit.String()
}

func (c Card) Title() string       { return zone.Mark(c.String(), c.String()) }
func (c Card) Description() string { return c.String() }

func (c Card) FilterValue() string { return zone.Mark(c.String(), c.Rank.String()) }
