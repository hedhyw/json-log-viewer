package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/app"
	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
	"github.com/hedhyw/json-log-viewer/internal/pkg/events"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
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

	fileName := ""
	var inputSource *source.Source

	switch flag.NArg() {
	case 0:
		// Tee stdin to a temp file, so that we can
		// lazy load the log entries using random access.
		fileName = "-"

		stdIn, err := getStdinReader(os.Stdin)
		if err != nil {
			fatalf("Stdin: %s\n", err)
		}

		inputSource, err = source.Reader(stdIn, cfg)
		if err != nil {
			fatalf("Could not create temp flie: %s\n", err)
		}
		defer inputSource.Close()

	case 1:
		fileName = flag.Arg(0)
		inputSource, err = source.File(fileName, cfg)
		if err != nil {
			fatalf("Could not create temp flie: %s\n", err)
		}
		defer inputSource.Close()

	default:
		fatalf("Invalid arguments, usage: %s file.log\n", os.Args[0])
	}

	appModel := app.NewModel(fileName, cfg, version)
	program := tea.NewProgram(appModel, tea.WithInputTTY(), tea.WithAltScreen())

	inputSource.StartStreaming(context.Background(), func(entries source.LazyLogEntries, err error) {
		if err != nil {
			program.Send(events.ErrorOccuredMsg{Err: err})
		} else {
			program.Send(events.LogEntriesUpdateMsg(entries))
		}
	})

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
