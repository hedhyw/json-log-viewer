package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/app"
	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source/fileinput"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source/readerinput"
)

const configFileName = ".jlv.jsonc"

func main() {
	configPath := flag.String("config", "", "Path to the config")
	flag.Parse()

	cfg, err := readConfig(*configPath)
	if err != nil {
		fatalf("Error reading config: %s\n", err)
	}

	var sourceInput source.Input

	switch flag.NArg() {
	case 0:
		sourceInput = readerinput.New(os.Stdin, cfg.StdinReadTimeout)
	case 1:
		sourceInput = fileinput.New(flag.Arg(0))
	default:
		fatalf("Invalid arguments, usage: %s file.log\n", os.Args[0])
	}

	appModel := app.NewModel(sourceInput, cfg)
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
func readConfig(configPath string) (*config.Config, error) {
	paths := []string{}

	if configPath != "" {
		paths = append(paths, configPath)
	}

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
