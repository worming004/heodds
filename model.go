package main

import (
	"errors"
	"iter"
	"log"
	"strconv"
	"strings"
	"time"

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
	errorMsg         string
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

func NewModel() Model {
	m := Model{}
	m.Cards = NewCards()

	return m
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

type ResetErrorMsg struct{}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	_, cmd = m.Cards.Update(msg)

	if m.Cards.HasChanged() {
		log.Println("card changed detected")
		selectedCard := m.Cards.selectedCard
		m.SetCard(selectedCard)
		m.ResetEquity()
	}

	switch msg := msg.(type) {
	case ResetErrorMsg:
		m.errorMsg = ""
	case tea.MouseMsg:
		if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft {
			log.Printf("Mouse coordinates: %v, %v\n", msg.X, msg.Y)
			if msg.Y == 4 && msg.X >= 21 && msg.X <= 25 {
				err := m.TriggerEvaluation()
				if err != nil {
					m.errorMsg = err.Error()
					return m, func() tea.Msg {
						time.Sleep(5 * time.Second)
						return ResetErrorMsg{}
					}
				}
			}
			if msg.Y == 5 {
				m.currentSelection = sp11
			}
			if msg.Y == 6 {
				m.currentSelection = sp12
			}
			if msg.Y == 8 {
				m.currentSelection = sp21
			}
			if msg.Y == 9 {
				m.currentSelection = sp22
			}

			if msg.Y == 12 {
				m.currentSelection = sf1
			}
			if msg.Y == 13 {
				m.currentSelection = sf2
			}
			if msg.Y == 14 {
				m.currentSelection = sf3
			}

			if msg.Y == 15 {
				m.currentSelection = st
			}
			if msg.Y == 16 {
				m.currentSelection = sr
			}
			if msg.Y == 17 {
				m.currentSelection = sp11
				m.ResetEquity()
				m.h1[0] = UnknownCard.c
				m.h1[1] = UnknownCard.c
				m.h2[0] = UnknownCard.c
				m.h2[1] = UnknownCard.c
				m.flop[0] = UnknownCard.c
				m.flop[1] = UnknownCard.c
				m.flop[2] = UnknownCard.c
				m.turn = UnknownCard.c
				m.river = UnknownCard.c
			}
		}
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if msg.String() == "enter" {
			err := m.TriggerEvaluation()
			if err != nil {
				m.errorMsg = err.Error()
				return m, func() tea.Msg {
					time.Sleep(5 * time.Second)
					return ResetErrorMsg{}
				}
			}
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

func (m *Model) TriggerEvaluation() error {
	err := m.EnsureNoDuplicate()
	if err != nil {
		return err
	}
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
		return err
	}
	m.equityj1 = res[0]
	m.equityj2 = res[1]

	return nil
}

func (m *Model) SetCard(c *Card) {
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
	writeEquity := func(v float64) {
		if v <= 0.2 {
			sb.WriteString("\x1b[38;5;52m") // dark red
		} else if v <= 0.4 {
			sb.WriteString("\x1b[38;5;88m") // red
		} else if v <= 0.48 {
			sb.WriteString("\x1b[38;5;130m") // orange
		} else if v <= 0.55 {
			sb.WriteString("\x1b[38;5;70m") // light green
		} else if v <= 0.65 {
			sb.WriteString("\x1b[38;5;76m") // green
		} else {
			sb.WriteString("\x1b[38;5;40m") // bright green
		}

		sb.WriteString(strconv.FormatFloat(v, 'f', 3, 64))

		sb.WriteString("\x1b[0m")
	}
	sb.WriteString(m.Cards.View())
	sb.WriteString("\n")

	sb.WriteString("currentSelection: ")
	currentSelectionString := strconv.Itoa(int(m.currentSelection))
	sb.WriteString(currentSelectionString)
	sb.WriteString("  |run|")
	sb.WriteString("\n")

	if m.currentSelection == sp11 {
		sb.WriteString("\x1b[33m")
	}
	sb.WriteString("J1 Hand 1: ")
	sb.WriteString(PokerCardString(m.h1[0]))
	sb.WriteString("\x1b[0m")
	sb.WriteString("\n")
	if m.currentSelection == sp12 {
		sb.WriteString("\x1b[33m")
	}
	sb.WriteString("J1 Hand 2: ")
	sb.WriteString(PokerCardString(m.h1[1]))
	sb.WriteString("\x1b[0m")
	sb.WriteString("\n")
	sb.WriteString("J1 equity: ")
	writeEquity(m.equityj1)
	sb.WriteString("\x1b[0m")
	sb.WriteString("\n")

	if m.currentSelection == sp21 {
		sb.WriteString("\x1b[33m")
	}
	sb.WriteString("J2 Hand 1: ")
	sb.WriteString(PokerCardString(m.h2[0]))
	sb.WriteString("\x1b[0m")
	sb.WriteString("\n")
	if m.currentSelection == sp22 {
		sb.WriteString("\x1b[33m")
	}
	sb.WriteString("J2 Hand 2: ")
	sb.WriteString(PokerCardString(m.h2[1]))
	sb.WriteString("\x1b[0m")
	sb.WriteString("\n")
	sb.WriteString("J2 equity: ")
	writeEquity(m.equityj2)
	sb.WriteString("\x1b[0m")
	sb.WriteString("\n")

	sb.WriteString("\n")

	if m.currentSelection == sf1 {
		sb.WriteString("\x1b[33m")
	}
	sb.WriteString("Flop 1: ")
	sb.WriteString(PokerCardString(m.flop[0]))
	sb.WriteString("\x1b[0m")
	sb.WriteString("\n")
	if m.currentSelection == sf2 {
		sb.WriteString("\x1b[33m")
	}
	sb.WriteString("Flop 2: ")
	sb.WriteString(PokerCardString(m.flop[1]))
	sb.WriteString("\x1b[0m")
	sb.WriteString("\n")
	if m.currentSelection == sf3 {
		sb.WriteString("\x1b[33m")
	}
	sb.WriteString("Flop 3: ")
	sb.WriteString(PokerCardString(m.flop[2]))
	sb.WriteString("\x1b[0m")
	sb.WriteString("\n")

	if m.currentSelection == st {
		sb.WriteString("\x1b[33m")
	}
	sb.WriteString("Turn: ")
	sb.WriteString(PokerCardString(m.turn))
	sb.WriteString("\x1b[0m")
	sb.WriteString("\n")

	if m.currentSelection == sr {
		sb.WriteString("\x1b[33m")
	}
	sb.WriteString("River: ")
	sb.WriteString(PokerCardString(m.river))
	sb.WriteString("\x1b[0m")
	sb.WriteString("\n")

	sb.WriteString("|Reset|\n")

	sb.WriteString("\x1b[38;5;88m") // red
	sb.WriteString(m.errorMsg)
	sb.WriteString("\x1b[0m")

	return sb.String()
}

func (m *Model) ResetEquity() {
	m.equityj1 = 0
	m.equityj2 = 0
}

func (m *Model) EnsureNoDuplicate() error {
	for i, c := range m.All() {
		for j, v := range m.All() {
			if j <= i {
				continue
			}
			if c == v && c != UnknownCard.c {
				return errors.New("duplicate card detected: " + PokerCardString(c) + PokerCardString(v))
			}
		}
	}

	return nil
}

func (m *Model) All() iter.Seq2[int, poker.Card] {
	return func(yield func(int, poker.Card) bool) {
		if !yield(0, m.h1[0]) {
			return
		}
		if !yield(1, m.h1[1]) {
			return
		}

		if !yield(2, m.h2[0]) {
			return
		}
		if !yield(3, m.h2[1]) {
			return
		}

		if !yield(4, m.flop[0]) {
			return
		}
		if !yield(5, m.flop[1]) {
			return
		}
		if !yield(6, m.flop[2]) {
			return
		}

		if !yield(7, m.turn) {
			return
		}

		if !yield(8, m.river) {
			return
		}
	}
}
