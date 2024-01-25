package main

import (
	"fmt"
	"os"
	"path"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/app"
	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
)

const configFileName = ".jlv.jsonc"

func main() {
	if len(os.Args) != 2 {
		fatalf("Invalid arguments, usage: %s file.log\n", os.Args[0])
	}

	cfg, err := readConfig()
	if err != nil {
		fatalf("Error reading config: %s\n", err)
	}

	appModel := app.NewModel(os.Args[1], cfg)
	program := tea.NewProgram(appModel, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		fatalf("Error running program: %s\n", err)
	}
}

func fatalf(message string, args ...any) {
	fmt.Fprintf(os.Stderr, message, args...)
	os.Exit(1)
}

// readConfig tries to read config from working directory or home directory.
// If configs are not found, then it returns a default configuration.
func readConfig() (*config.Config, error) {
	paths := []string{}

	workDir, err := os.Getwd()
	if err == nil {
		paths = append(paths, path.Join(workDir, configFileName))
	}

	homeDir, err := os.UserHomeDir()
	if err == nil {
		paths = append(paths, path.Join(homeDir, configFileName))
	}

	return config.Read(paths...)
}
