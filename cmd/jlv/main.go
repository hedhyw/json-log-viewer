package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
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

	err := runApp(applicationArguments{
		Stdout: os.Stdout,
		Stdin:  os.Stdin,

		ConfigPath:   *configPath,
		PrintVersion: *printVersion,
		Args:         flag.Args(),

		RunProgram: func(p *tea.Program) (tea.Model, error) {
			return p.Run()
		},
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: "+err.Error())
		os.Exit(1)
	}
}

type applicationArguments struct {
	Stdout io.Writer
	Stdin  fs.File

	ConfigPath   string
	PrintVersion bool
	Args         []string

	RunProgram func(*tea.Program) (tea.Model, error)
}

func runApp(args applicationArguments) (err error) {
	if args.PrintVersion {
		// nolint: forbidigo // Version command.
		fmt.Fprintln(args.Stdout, "github.com/hedhyw/json-log-viewer@"+version)

		return nil
	}

	cfg, err := readConfig(args.ConfigPath)
	if err != nil {
		return fmt.Errorf("reading config: %w", err)
	}

	fileName := ""
	var inputSource *source.Source

	switch len(args.Args) {
	case 0:
		// Tee stdin to a temp file, so that we can
		// lazy load the log entries using random access.
		fileName = "-"

		stdin, err := getStdinReader(args.Stdin)
		if err != nil {
			return fmt.Errorf("getting stdin: %w", err)
		}

		inputSource, err = source.Reader(stdin, cfg)
		if err != nil {
			return fmt.Errorf("creating a temporary file: %w", err)
		}

		defer func() { err = errors.Join(err, inputSource.Close()) }()
	case 1:
		fileName = args.Args[0]

		inputSource, err = source.File(fileName, cfg)
		if err != nil {
			return fmt.Errorf("reading file: %w", err)
		}

		defer func() { err = errors.Join(err, inputSource.Close()) }()
	default:
		// nolint: err113 // One time case.
		return fmt.Errorf("invalid arguments, usage: %s file.log", os.Args[0])
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

	if _, err := args.RunProgram(program); err != nil {
		return fmt.Errorf("running program: %w", err)
	}

	return nil
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
