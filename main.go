package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	zone "github.com/lrstanley/bubblezone"
)

func init() {
	zone.NewGlobal()
}

func main() {
	m := NewModel()

	if _, err := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run(); err != nil {
		fmt.Println("error running program:", err)
		os.Exit(1)
	}
}
