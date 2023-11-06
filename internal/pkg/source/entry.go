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
		return formatMessage(formatSecondTime(value))
	case config.FieldKindMilliTime:
		return formatMessage(formatMillisecondTime(value))
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

func formatSecondTime(timeStr string) string {
	// Parse the string to float64
	seconds, err := strconv.ParseFloat(timeStr, 64)
	if err != nil {
		log.Println("error parsing floating: " + err.Error())
		panic("Could not parse float")
	}

	// Separate the seconds into integer and fractional parts
	secInt := int64(seconds)
	secFrac := seconds - float64(secInt)

	// Convert fractional seconds to nanoseconds
	nanoSec := int64(secFrac * float64(time.Second))

	time := time.Unix(secInt, nanoSec)
	log.Println("input: " + timeStr)
	log.Println("seconds: " + fmt.Sprint(seconds))
	log.Println("nanoSec: " + fmt.Sprint(nanoSec))

	return toRfc3339(time)
}

func formatMillisecondTime(timeStr string) string {
	// This could be seconds or it could be a float
	milliseconds, err := strconv.ParseFloat(timeStr, 64)
	if err != nil {
		return timeStr
	}

	seconds := int64(milliseconds / 1000)

	nanoseconds := int64((milliseconds - float64(seconds*1000)) * 1e6) // leftover milliseconds to nanoseconds

	time := time.Unix(seconds, nanoseconds)
	log.Println("input: " + timeStr)
	log.Println("seconds: " + fmt.Sprint(seconds))
	log.Println("milliseconds: " + fmt.Sprint(milliseconds))
	log.Println("Nanoseconds: " + fmt.Sprint(nanoseconds))

	return toRfc3339(time)
}

func toRfc3339(t time.Time) string {
	return t.Format(time.RFC3339)
}
