package main

import (
	"bytes"
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

// version will be set on build.
var version = "development"

const configFileName = ".jlv.jsonc"

func main() {
	configPath := flag.String("config", "", "Path to the config")
	printVersion := flag.Bool("version", false, "Print version")
	flag.Parse()

	if *printVersion {
		// nolint: forbidigo // Version command.
		print("github.com/hedhyw/json-log-viewer@" + version + "\n")

		return
	}

	cfg, err := readConfig(*configPath)
	if err != nil {
		fatalf("Error reading config: %s\n", err)
	}

	var sourceInput source.Input

	switch flag.NArg() {
	case 0:
		sourceInput, err = getStdinSource(cfg)
		if err != nil {
			fatalf("Stdin: %s\n", err)
		}
	case 1:
		sourceInput = fileinput.New(flag.Arg(0))
	default:
		fatalf("Invalid arguments, usage: %s file.log\n", os.Args[0])
	}

	appModel := app.NewModel(sourceInput, cfg, version)
	program := tea.NewProgram(appModel, tea.WithInputTTY(), tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		fatalf("Error running program: %s\n", err)
	}
}

func getStdinSource(cfg *config.Config) (source.Input, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat: %w", err)
	}

	if stat.Mode()&os.ModeNamedPipe == 0 {
		return readerinput.New(bytes.NewReader(nil), cfg.StdinReadTimeout), nil
	}

	return readerinput.New(os.Stdin, cfg.StdinReadTimeout), nil
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
