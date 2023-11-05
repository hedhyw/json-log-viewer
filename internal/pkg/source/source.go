package source

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
)

const (
	maxLineSize = 8 * 1024 * 1024

	logEntriesEstimateNumber = 256
)

// LoadLogsFromFile loads json log entries from file.
func LoadLogsFromFile(
	path string,
	cfg *config.Config,
) func() tea.Msg {
	return func() (msg tea.Msg) {
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("os: %w", err)
		}

		defer file.Close()

		logEntries, err := parseLogEntriesFromReader(file, cfg)
		if err != nil {
			return fmt.Errorf("parsing from reader: %w", err)
		}

		return logEntries.Reverse()
	}
}

func parseLogEntriesFromReader(
	reader io.Reader,
	cfg *config.Config,
) (LogEntries, error) {
	bufReader := bufio.NewReaderSize(reader, maxLineSize)
	logEntries := make(LogEntries, 0, logEntriesEstimateNumber)

	for {
		line, _, err := bufReader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, fmt.Errorf("reading line: %w", err)
		}

		if len(bytes.TrimSpace(line)) > 0 {
			logEntries = append(logEntries, ParseLogEntry(line, cfg))
		}
	}

	return logEntries, nil
}
