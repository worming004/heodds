package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func setLoggerToFile() {
	f, err := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("error opening file:", err)
		return
	}

	// Set the output of log to the file
	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Logging to debug.log")
}

func main() {
	m := NewModel()
	setLoggerToFile()

	if _, err := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run(); err != nil {
		fmt.Println("error running program:", err)
		os.Exit(1)
	}
}
