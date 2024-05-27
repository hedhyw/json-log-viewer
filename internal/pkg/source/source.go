package source

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

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
) (_ LazyLogEntries, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("os: %w", err)
	}

	defer func() { err = errors.Join(err, file.Close()) }()

	logEntries, err := ParseLogEntriesFromReader(file, cfg)
	if err != nil {
		return nil, fmt.Errorf("parsing from reader: %w", err)
	}

	return logEntries.Reverse(), nil
}

func ParseLogEntriesFromReader(
	reader io.Reader,
	cfg *config.Config,
) (LazyLogEntries, error) {
	reader = io.LimitReader(reader, cfg.MaxFileSizeBytes)

	bufReader := bufio.NewReaderSize(reader, maxLineSize)
	logEntries := make(LazyLogEntries, 0, logEntriesEstimateNumber)

	for {
		line, _, err := bufReader.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, fmt.Errorf("reading line: %w", err)
		}

		line = bytes.TrimSpace(line)

		if len(line) > 0 {
			lineClone := make([]byte, len(line))
			copy(lineClone, line)

			logEntries = append(logEntries, LazyLogEntry{Line: lineClone})
		}
	}

	return logEntries, nil
}
