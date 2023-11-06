package source

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/charmbracelet/bubbles/table"
	"github.com/yalp/jsonpath"

	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
)

// LogEntry is a single partly-parse record of the log.
type LogEntry struct {
	Fields []string
	Line   json.RawMessage
}

// Row returns table.Row representation of the log entry.
func (e LogEntry) Row() table.Row {
	return e.Fields
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

func parseField(parsedLine any, field config.Field) string {
	for _, ref := range field.References {
		foundField, err := jsonpath.Read(parsedLine, ref)
		if err != nil {
			continue
		}

		jsonField, err := json.Marshal(foundField)
		if err != nil {
			return fmt.Sprint(field)
		}

		unquotedField, err := strconv.Unquote(string(jsonField))
		// It's possible that what were were given is an integer or float
		// in which case, calling Unquote isn't doing us a lot of good.
		// Therefore, we just convert to a string value and proceed.
		if err != nil {
			unquotedField = string(jsonField)
		}

		return formatField(unquotedField, field.Kind)
	}

	return "-"
}

func formatField(
	value string,
	kind config.FieldKind,
) string {
	value = strings.TrimSpace(value)

	switch kind {
	case config.FieldKindMessage:
		return formatMessage(value)
	case config.FieldKindLevel:
		return string(ParseLevel(formatMessage(value)))
	case config.FieldKindTime:
		return formatMessage(value)
	case config.FieldKindSecondTime:
		return formatMessage(formatTimeString(value, "s"))
	case config.FieldKindMilliTime:
		return formatMessage(formatTimeString(value, "ms"))
	case config.FieldKindMicroTime:
		return formatMessage(formatTimeString(value, "us"))
	case config.FieldKindAny:
		return formatMessage(value)
	default:
		return formatMessage(value)
	}
}

// ParseLogEntry parses a single log entry from the json line.
func ParseLogEntry(
	line json.RawMessage,
	cfg *config.Config,
) LogEntry {
	var parsedLine any

	err := json.Unmarshal(normalizeJSON(line), &parsedLine)
	if err != nil {
		return getPlainLogEntry(line, cfg)
	}

	fields := make([]string, 0, len(cfg.Fields))

	for _, f := range cfg.Fields {
		fields = append(fields, parseField(parsedLine, f))
	}

	return LogEntry{
		Line:   line,
		Fields: fields,
	}
}

func getPlainLogEntry(
	line json.RawMessage,
	cfg *config.Config,
) LogEntry {
	fields := make([]string, len(cfg.Fields))

	for i, f := range cfg.Fields {
		fields[i] = "-"

		if f.Kind == config.FieldKindMessage {
			fields[i] = string(line)
		}
	}

	return LogEntry{
		Fields: fields,
		Line:   line,
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

func formatTimeString(timeStr string, unit string) string {
	duration, err := time.ParseDuration(timeStr + unit)
	if err != nil {
		log.Println("Error parsing time: " + timeStr + " unit: " + unit + "\nError: " + err.Error())
		return timeStr
	}

	seconds := int64(duration.Seconds())
	nanoseconds := duration.Nanoseconds()

	var t time.Time = time.Unix(seconds, nanoseconds-(seconds*int64(time.Second)))

	return t.Format(time.RFC3339)
}
