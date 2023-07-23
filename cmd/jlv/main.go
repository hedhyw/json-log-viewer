package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/app"
)

func main() {
	if len(os.Args) != 2 {
		fatalf("Invalid arguments, usage: %s file.log\n", os.Args[0])
	}

	appModel := app.NewModel(os.Args[1])
	program := tea.NewProgram(appModel, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		fatalf("Error running program: %s\n", err)
	}
}

func fatalf(message string, args ...any) {
	fmt.Fprintf(os.Stderr, message, args...)
	os.Exit(1)
}
