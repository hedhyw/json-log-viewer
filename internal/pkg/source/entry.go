package source

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/charmbracelet/bubbles/table"
	"github.com/yalp/jsonpath"

	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
)

const (
	unitSeconds = "s"
	unitMilli   = "ms"
	unitMicro   = "us"
)

// LazyLogEntry holds unredenred LogEntry. Use `LogEntry` getter.
type LazyLogEntry struct {
	offset int64
	length int
	index  int
}

// Length of the entry.
func (e LazyLogEntry) Length() int {
	return e.length
}

// Line re-reads the line.
func (e LazyLogEntry) Line(file *os.File) (json.RawMessage, error) {
	data := make([]byte, e.length)

	_, err := file.ReadAt(data, e.offset)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// LogEntry parses and returns `LogEntry`.
func (e LazyLogEntry) LogEntry(file *os.File, cfg *config.Config) LogEntry {
	line, err := e.Line(file)
	if err != nil {
		return LogEntry{
			Index: e.index,
			Error: err,
		}
	}

	entry := parseLogEntry(line, cfg)
	entry.Index = e.index

	return entry
}

// LogEntry is a single partly-parse record of the log.
type LogEntry struct {
	Index  int
	Fields []string
	Line   json.RawMessage
	Error  error
}

// Row returns table.Row representation of the log entry.
func (e LogEntry) Row() table.Row {
	return e.Fields
}

// LazyLogEntries is a helper type definition for the slice of lazy log entries.
type LazyLogEntries struct {
	Seeker  *os.File
	Entries []LazyLogEntry
}

// Row returns table.Row representation of the log entry.
func (entries LazyLogEntries) Row(cfg *config.Config, i int) table.Row {
	return entries.LogEntry(cfg, i).Fields
}

// LogEntry getter.
func (entries LazyLogEntries) LogEntry(cfg *config.Config, i int) LogEntry {
	return entries.Entries[i].LogEntry(entries.Seeker, cfg)
}

// Len return the number of all entries.
func (entries LazyLogEntries) Len() int {
	return len(entries.Entries)
}

// Filter filters entries by a case-insensitive substring match. A term
// wrapped in slashes (/.../) is matched as a regular expression instead.
//
// In fulltext mode the term is matched against the whole raw JSON line, in
// field mode against the rendered value of the given field.
func (entries LazyLogEntries) Filter(term string, fieldName string, c *config.Config) (LazyLogEntries, error) {
	if term == "" {
		return entries, nil
	}

	fieldIndex := getFilterFieldNameIndex(fieldName, c)
	if len(fieldName) != 0 && fieldIndex < 0 {
		return LazyLogEntries{}, fmt.Errorf("%w: unknown field: %s", ErrInvalidFilter, fieldName)
	}

	matches, err := NewMatcher(term)
	if err != nil {
		return LazyLogEntries{}, err
	}

	filtered := make([]LazyLogEntry, 0, len(entries.Entries))

	for _, f := range entries.Entries {
		var value []byte

		if len(fieldName) == 0 {
			// Fulltext mode. The stored line keeps the trailing line break,
			// it is trimmed so that the `$` anchor can match the end of it.
			line, err := f.Line(entries.Seeker)
			if err != nil {
				return LazyLogEntries{}, err
			}

			value = bytes.TrimRight(line, "\r\n")
		} else {
			// Field mode.
			entry := f.LogEntry(entries.Seeker, c)
			if entry.Error != nil {
				return LazyLogEntries{}, entry.Error
			}

			// A field of a non-JSON line holds the untouched line, which
			// still keeps its trailing line break, so it is trimmed too.
			value = bytes.TrimRight([]byte(entry.Fields[fieldIndex]), "\r\n")
		}

		if matches(value) {
			filtered = append(filtered, f)
		}
	}

	return LazyLogEntries{
		Seeker:  entries.Seeker,
		Entries: filtered,
	}, nil
}

// NewMatcher returns a case-insensitive predicate for the given term. A term
// wrapped in slashes (/.../) is compiled as a regular expression, otherwise
// a substring match is used.
//
// It returns an error wrapping ErrInvalidFilter if the term holds a malformed
// regular expression. Callers can validate a term up front to report the
// problem while the user can still correct it.
func NewMatcher(term string) (func(value []byte) bool, error) {
	// A term is a regular expression only if it is wrapped in slashes and
	// holds at least one character in between, so that `/` and `//` stay
	// plain substring searches.
	if len(term) > 2 && term[0] == '/' && term[len(term)-1] == '/' {
		expr := term[1 : len(term)-1]

		exp, err := regexp.Compile("(?i)" + expr)
		if err != nil {
			// Report the failure against the expression that the user typed,
			// without the case-insensitivity flag that is injected above.
			if _, userErr := regexp.Compile(expr); userErr != nil {
				err = userErr
			}

			return nil, fmt.Errorf("%w: compiling regular expression: %w", ErrInvalidFilter, err)
		}

		return exp.Match, nil
	}

	termLower := bytes.ToLower([]byte(term))

	return func(value []byte) bool {
		return bytes.Contains(bytes.ToLower(value), termLower)
	}, nil
}

func getFilterFieldNameIndex(fieldName string, c *config.Config) int {
	if c == nil || fieldName == "" {
		return -1
	}

	for i, field := range c.Fields {
		if strings.EqualFold(field.Title, fieldName) {
			return i
		}
	}

	return -1
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

		return formatField(unquotedField, field, cfg)
	}

	return "-"
}

//nolint:cyclop // The cyclomatic complexity here is so high because of the number of FieldKinds.
func formatField(
	value string,
	field config.Field,
	cfg *config.Config,
) string {
	kind := field.Kind
	value = strings.TrimSpace(value)

	timeFormat := config.DefaultTimeFormat

	if field.TimeFormat != nil {
		timeFormat = *field.TimeFormat
	}

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
		return formatMessage(reformatTime(value, cfg.TimeLayouts, timeFormat))
	case config.FieldKindSecondTime:
		return formatMessage(formatTimeValue(value, unitSeconds, timeFormat))
	case config.FieldKindMilliTime:
		return formatMessage(formatTimeValue(value, unitMilli, timeFormat))
	case config.FieldKindMicroTime:
		return formatMessage(formatTimeValue(value, unitMicro, timeFormat))
	case config.FieldKindAny:
		return formatMessage(value)
	default:
		return formatMessage(value)
	}
}

func reformatTime(value string, layoutsToReformat []string, timeFormat string) string {
	for _, laoyout := range layoutsToReformat {
		parsed, err := time.Parse(laoyout, value)
		if err == nil {
			return parsed.Format(timeFormat)
		}
	}

	return value
}

// parseLogEntry parses a single log entry from the json line.
func parseLogEntry(
	line json.RawMessage,
	cfg *config.Config,
) LogEntry {
	var parsedLine any

	err := json.Unmarshal(normalizeJSON(line), &parsedLine)
	if err != nil {
		return getPlainLogEntry(line, cfg)
	}

	if _, ok := parsedLine.(map[string]any); !ok {
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

	const (
		unixSecondsLength = 10
		unixMilliLength   = 13
		unixMicroLength   = 16
	)

	switch {
	case intLength <= unixSecondsLength:
		return config.FieldKindSecondTime
	case intLength > unixSecondsLength && intLength <= unixMilliLength:
		return config.FieldKindMilliTime
	case intLength > unixMilliLength && intLength <= unixMicroLength:
		return config.FieldKindMicroTime
	default:
		return config.FieldKindTime
	}
}

func formatTimeValue(timeValue string, unit string, format string) string {
	duration, err := time.ParseDuration(timeValue + unit)
	if err != nil {
		return timeValue
	}

	return time.UnixMilli(0).Add(duration).UTC().Format(format)
}
