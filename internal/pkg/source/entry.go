package source

import (
	"bytes"
	"encoding/json"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/table"
	"github.com/valyala/fastjson"
)

// LogEntry is a single partly-parse record of the log.
type LogEntry struct {
	Time    string
	Level   Level
	Message string
	Line    json.RawMessage
}

// Row returns table.Row representation of the log entry.
func (e LogEntry) Row() table.Row {
	return table.Row{
		e.Time,
		string(e.Level),
		e.Message,
	}
}

// LogEntries is a helper type definition for the slice of log entries.
type LogEntries []LogEntry

// Filter filters entries by ignore case exact match.
func (entries LogEntries) Filter(term string) LogEntries {
	if term == "" {
		return entries
	}

	termLower := bytes.ToLower([]byte(term))

	filtered := make([]LogEntry, 0, len(entries))

	for _, f := range entries {
		if bytes.Contains(bytes.ToLower(f.Line), termLower) {
			filtered = append(filtered, f)
		}
	}

	return filtered
}

func (entries LogEntries) Reverse() LogEntries {
	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
	}

	return entries
}

// Rows returns all table.Row by log entries.
func (entries LogEntries) Rows() []table.Row {
	rows := make([]table.Row, len(entries))

	for i, e := range entries {
		rows[i] = e.Row()
	}

	return rows
}

// ParseLogEntry parses a single log entry from the json line.
func ParseLogEntry(line json.RawMessage) LogEntry {
	var jsonParser fastjson.Parser

	lineToParse := make([]byte, len(line))
	copy(lineToParse, line)
	line = lineToParse

	value, err := jsonParser.ParseBytes(lineToParse)
	if err != nil {
		return LogEntry{
			Line:    line,
			Time:    "-",
			Message: formatMessage(string(line)),
			Level:   LevelUnknown,
		}
	}

	return LogEntry{
		Line:    line,
		Time:    formatMessage(extractTime(value)),
		Message: formatMessage(extractMessage(value)),
		Level:   extractLevel(value),
	}
}

func formatMessage(msg string) string {
	msg = strings.NewReplacer("\n", "\\n", "\t", "\\t").Replace(msg)

	msg = strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}

		return -1
	}, msg)

	return msg
}
