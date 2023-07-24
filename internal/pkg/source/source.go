package source

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// LoadLogsFromFile loads json log entries from file.
func LoadLogsFromFile(path string) func() tea.Msg {
	return func() (msg tea.Msg) {
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("os: %w", err)
		}

		defer file.Close()

		scanner := bufio.NewScanner(file)

		logEntries := make(LogEntries, 0, 256)

		for scanner.Scan() {
			line := scanner.Bytes()

			if len(bytes.TrimSpace(line)) > 0 {
				logEntries = append(logEntries, ParseLogEntry(line))
			}
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("scanning: %w", err)
		}

		return logEntries.Reverse()
	}
}
