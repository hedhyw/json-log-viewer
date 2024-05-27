package source

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/charmbracelet/bubbles/table"
	"github.com/yalp/jsonpath"

	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
)

// LazyLogEntry holds unredenred LogEntry. Use `LogEntry` getter.
type LazyLogEntry struct {
	Line json.RawMessage
}

// LogEntry parses and returns `LogEntry`.
func (e LazyLogEntry) LogEntry(cfg *config.Config) LogEntry {
	return ParseLogEntry(e.Line, cfg)
}

// Row returns table.Row representation of the log entry.
func (e LazyLogEntry) Row(cfg *config.Config) table.Row {
	return e.LogEntry(cfg).Fields
}

// LogEntry is a single partly-parse record of the log.
type LogEntry struct {
	Fields []string
	Line   json.RawMessage
}

// Row returns table.Row representation of the log entry.
func (e LogEntry) Row() table.Row {
	return e.Fields
}

// LazyLogEntries is a helper type definition for the slice of lazy log entries.
type LazyLogEntries []LazyLogEntry

// Filter filters entries by ignore case exact match.
func (entries LazyLogEntries) Filter(term string) LazyLogEntries {
	if term == "" {
		return entries
	}

	termLower := bytes.ToLower([]byte(term))

	filtered := make(LazyLogEntries, 0, len(entries))

	for _, f := range entries {
		if bytes.Contains(bytes.ToLower(f.Line), termLower) {
			filtered = append(filtered, f)
		}
	}

	return filtered
}

// Reverse all entries.
func (entries LazyLogEntries) Reverse() LazyLogEntries {
	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
	}

	return entries
}

func parseField(
	parsedLine any,
	field config.Field,
	cfg *config.Config,
) string {
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
		// It's possible that what we were given is an integer or float
		// in which case, calling Unquote isn't doing us a lot of good.
		// Therefore, we just convert to a string value and proceed.
		if err != nil {
			unquotedField = string(jsonField)
		}

		return formatField(unquotedField, field.Kind, cfg)
	}

	return "-"
}

//nolint:cyclop // The cyclomatic complexity here is so high because of the number of FieldKinds.
func formatField(
	value string,
	kind config.FieldKind,
	cfg *config.Config,
) string {
	value = strings.TrimSpace(value)

	// Numeric time attempts to infer the duration based on the length of the string
	if kind == config.FieldKindNumericTime {
		kind = guessTimeFieldKind(value)
	}

	switch kind {
	case config.FieldKindMessage:
		return formatMessage(value)
	case config.FieldKindLevel:
		return string(ParseLevel(formatMessage(value), cfg.CustomLevelMapping))
	case config.FieldKindTime:
		return formatMessage(value)
	case config.FieldKindNumericTime:
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
		fields = append(fields, parseField(parsedLine, f, cfg))
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

// We can only guess the time via a heuristic. We do this by looking at the number of digits
// (before the decimal point) in the string. This is far from perfect.
func guessTimeFieldKind(timeStr string) config.FieldKind {
	intValue, err := strconv.ParseInt(strings.Split(timeStr, ".")[0], 10, 64)
	if err != nil {
		return config.FieldKindTime
	}

	if intValue <= 0 {
		return config.FieldKindTime
	}

	intLength := len(strconv.FormatInt(intValue, 10))

	switch {
	case intLength <= 10:
		return config.FieldKindSecondTime
	case intLength > 10 && intLength <= 13:
		return config.FieldKindMilliTime
	case intLength > 13 && intLength <= 16:
		return config.FieldKindMicroTime
	default:
		return config.FieldKindTime
	}
}

func formatTimeString(timeStr string, unit string) string {
	duration, err := time.ParseDuration(timeStr + unit)
	if err != nil {
		return timeStr
	}

	return time.UnixMilli(0).Add(duration).UTC().Format(time.RFC3339)
}
