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
	"strings"

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
		// Multiple files are concatenated and teed to a temporary file,
		// so that we can lazy load the log entries using random access.
		fileName = fmt.Sprintf("%s (+%d)", args.Args[0], len(args.Args)-1)

		reader, closeFiles, errOpen := openLogFiles(args.Args)
		if errOpen != nil {
			return fmt.Errorf("reading files: %w", errOpen)
		}

		defer func() { err = errors.Join(err, closeFiles()) }()

		inputSource, err = source.Reader(reader, cfg)
		if err != nil {
			return fmt.Errorf("creating a temporary file: %w", err)
		}

		defer func() { err = errors.Join(err, inputSource.Close()) }()
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

// openLogFiles opens all given files and returns a reader that reads them
// one after another. The files are separated by a line break, because the
// last line of a file is not guaranteed to have one. Empty lines are
// skipped while parsing.
//
// The returned closer closes all opened files.
func openLogFiles(names []string) (io.Reader, func() error, error) {
	files := make([]*os.File, 0, len(names))

	closeFiles := func() error {
		errMulti := make([]error, 0, len(files))

		for _, f := range files {
			errMulti = append(errMulti, f.Close())
		}

		return errors.Join(errMulti...)
	}

	readers := make([]io.Reader, 0, 2*len(names))

	for _, name := range names {
		file, err := os.Open(name)
		if err != nil {
			return nil, nil, errors.Join(fmt.Errorf("opening: %w", err), closeFiles())
		}

		files = append(files, file)
		readers = append(readers, file, strings.NewReader("\n"))
	}

	return io.MultiReader(readers...), closeFiles, nil
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
