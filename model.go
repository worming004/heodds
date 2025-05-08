package main

import (
	"log"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/whywaita/poker-go"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type Model struct {
	Cards   *Cards
	players []poker.Player

	h1 [2]poker.Card
	h2 [2]poker.Card

	equityj1 float64
	equityj2 float64

	flop  [3]poker.Card
	turn  poker.Card
	river poker.Card

	currentSelection selection
}

func IsUnknownCard(c poker.Card) bool {
	return c.Rank == poker.RankUnknown
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

func (m *Model) SetNextSelection() {
	if m.currentSelection == 8 {
		m.currentSelection = sp11
		return
	}

	m.currentSelection = m.currentSelection + 1
}
func (m *Model) SetPreviousSelection() {
	if m.currentSelection == 0 {
		m.currentSelection = sr
		return
	}

	m.currentSelection = m.currentSelection - 1
}

func removeCards(d *poker.Deck, count int) {
	for range count {
		d.DrawCard()
	}
}

func NewModel() Model {
	m := Model{}
	m.Cards = NewCards()

	return m
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	_, cmd = m.Cards.Update(msg)

	if m.Cards.HasChanged() {
		log.Println("card changed detected")
		selectedCard := m.Cards.selectedCard
		m.setCard(selectedCard)
		m.resetEquity()
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if msg.String() == "enter" {
			p1 := *poker.NewPlayer("Player 1", m.h1[:])
			p2 := *poker.NewPlayer("Player 2", m.h2[:])

			community := make([]poker.Card, 0, 5)
			if !IsUnknownCard(m.flop[0]) {
				community = append(community, m.flop[:]...)
			}

			if !IsUnknownCard(m.turn) {
				community = append(community, m.turn)
			}
			if !IsUnknownCard(m.river) {
				community = append(community, m.river)
			}

			res, err := poker.EvaluateEquityByMadeHandWithCommunity([]poker.Player{p1, p2}, community)
			if err != nil {
				panic(err)
			}
			m.equityj1 = res[0]
			m.equityj2 = res[1]
		}
		if msg.String() == "left" {
			m.SetPreviousSelection()
		}
		if msg.String() == "right" {
			m.SetNextSelection()
		}
	}

	return m, cmd
}

func (m *Model) setCard(c *Card) {
	log.Println("setCard for selection:", c, m.currentSelection)
	switch m.currentSelection {
	case sp11:
		m.h1[0] = c.c
	case sp12:
		m.h1[1] = c.c
	case sp21:
		m.h2[0] = c.c
	case sp22:
		m.h2[1] = c.c
	case sf1:
		m.flop[0] = c.c
	case sf2:
		m.flop[1] = c.c
	case sf3:
		m.flop[2] = c.c
	case st:
		m.turn = c.c
	case sr:
		m.river = c.c
	}

	m.SetNextSelection()
}

// View implements tea.Model.
func (m Model) View() string {
	sb := strings.Builder{}
	sb.WriteString(m.Cards.View())
	sb.WriteString("\n")

	sb.WriteString("currentSelection: ")
	currentSelectionString := strconv.Itoa(int(m.currentSelection))
	sb.WriteString(currentSelectionString)
	sb.WriteString("\n")

	sb.WriteString("J1 Hand 1: ")
	sb.WriteString(PokerCardString(m.h1[0]))
	sb.WriteString("\n")
	sb.WriteString("J1 Hand 2: ")
	sb.WriteString(PokerCardString(m.h1[1]))
	sb.WriteString("\n")
	sb.WriteString("J1 equity: ")
	sb.WriteString(strconv.FormatFloat(m.equityj1, 'f', 3, 64))
	sb.WriteString("\n")

	sb.WriteString("J2 Hand 1: ")
	sb.WriteString(PokerCardString(m.h2[0]))
	sb.WriteString("\n")
	sb.WriteString("J2 Hand 2: ")
	sb.WriteString(PokerCardString(m.h2[1]))
	sb.WriteString("\n")
	sb.WriteString("J2 equity: ")
	sb.WriteString(strconv.FormatFloat(m.equityj2, 'f', 3, 64))
	sb.WriteString("\n")

	sb.WriteString("\n")

	sb.WriteString("Flop 1: ")
	sb.WriteString(PokerCardString(m.flop[0]))
	sb.WriteString("\n")
	sb.WriteString("Flop 2: ")
	sb.WriteString(PokerCardString(m.flop[1]))
	sb.WriteString("\n")
	sb.WriteString("Flop 3: ")
	sb.WriteString(PokerCardString(m.flop[2]))
	sb.WriteString("\n")

	sb.WriteString("Turn: ")
	sb.WriteString(PokerCardString(m.turn))
	sb.WriteString("\n")

	sb.WriteString("River: ")
	sb.WriteString(PokerCardString(m.river))
	sb.WriteString("\n")

	return sb.String()
}

func (m *Model) resetEquity() {
	m.equityj1 = 0
	m.equityj2 = 0
}
