package source_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/hedhyw/json-log-viewer/internal/pkg/source"

	"github.com/stretchr/testify/assert"
)

func TestParseLogEntry(t *testing.T) {
	t.Parallel()

	testCases := [...]struct {
		Name   string
		JSON   string
		Assert func(tb testing.TB, logEntry source.LogEntry)
	}{{
		Name: "plain_log",
		JSON: "Hello World",
		Assert: func(tb testing.TB, logEntry source.LogEntry) {
			tb.Helper()

			assert.Equal(t, "Hello World", logEntry.Message)
			assert.Equal(t, source.LevelUnknown, logEntry.Level)
			assert.Equal(t, "-", logEntry.Time)
		},
	}, {
		Name: "time_number",
		JSON: `{"time":1}`,
		Assert: func(tb testing.TB, logEntry source.LogEntry) {
			tb.Helper()

			assert.Equal(t, "1", logEntry.Time)
		},
	}, {
		Name: "timestamp_number",
		JSON: `{"timestamp":1}`,
		Assert: func(tb testing.TB, logEntry source.LogEntry) {
			tb.Helper()

			assert.Equal(t, "1", logEntry.Time)
		},
	}, {
		Name: "time_text",
		JSON: `{"time":"1970-01-01T00:00:00.00"}`,
		Assert: func(tb testing.TB, logEntry source.LogEntry) {
			tb.Helper()

			assert.Equal(t, "1970-01-01T00:00:00.00", logEntry.Time)
		},
	}, {
		Name: "timestamp_text",
		JSON: `{"timestamp":"1970-01-01T00:00:00.00"}`,
		Assert: func(tb testing.TB, logEntry source.LogEntry) {
			tb.Helper()

			assert.Equal(t, "1970-01-01T00:00:00.00", logEntry.Time)
		},
	}, {
		Name: "message",
		JSON: `{"message":"message"}`,
		Assert: func(tb testing.TB, logEntry source.LogEntry) {
			tb.Helper()

			assert.Equal(t, "message", logEntry.Message)
		},
	}, {
		Name: "msg",
		JSON: `{"msg":"msg"}`,
		Assert: func(tb testing.TB, logEntry source.LogEntry) {
			tb.Helper()

			assert.Equal(t, "msg", logEntry.Message)
		},
	}, {
		Name: "message_special_rune",
		JSON: `{"message":"mes` + string(rune(1)) + `sage"}`,
		Assert: func(tb testing.TB, logEntry source.LogEntry) {
			tb.Helper()

			assert.Equal(t, "message", logEntry.Message)
		},
	}, {
		Name: "error",
		JSON: `{"error":"error"}`,
		Assert: func(tb testing.TB, logEntry source.LogEntry) {
			tb.Helper()

			assert.Equal(t, "error", logEntry.Message)
		},
	}, {
		Name: "err",
		JSON: `{"err":"err"}`,
		Assert: func(tb testing.TB, logEntry source.LogEntry) {
			tb.Helper()

			assert.Equal(t, "err", logEntry.Message)
		},
	}, {
		Name: "level",
		JSON: `{"level":"INFO"}`,
		Assert: func(tb testing.TB, logEntry source.LogEntry) {
			tb.Helper()

			assert.Equal(t, "info", logEntry.Level.String())
		},
	}}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			actual := source.ParseLogEntry(json.RawMessage(testCase.JSON))
			testCase.Assert(t, actual)
		})
	}
}

func TestLogEntryRow(t *testing.T) {
	t.Parallel()

	entry := getFakeLogEntry()
	row := entry.Row()

	if assert.Len(t, row, 3) {
		assert.Equal(t, entry.Time, row[0])
		assert.Equal(t, string(entry.Level), row[1])
		assert.Equal(t, entry.Message, row[2])
	}
}

func TestLogEntriesRows(t *testing.T) {
	t.Parallel()

	entries := source.LogEntries{
		getFakeLogEntry(),
		getFakeLogEntry(),
		getFakeLogEntry(),
	}
	rows := entries.Rows()

	if assert.Len(t, rows, len(entries)) {
		for i, e := range entries {
			assert.Equal(t, e.Row(), rows[i])
		}
	}
}

func TestLogEntriesReverse(t *testing.T) {
	t.Parallel()

	t.Run("simple", func(t *testing.T) {
		t.Parallel()

		original := source.LogEntries{
			getFakeLogEntry(),
			getFakeLogEntry(),
			getFakeLogEntry(),
		}

		entries := make(source.LogEntries, len(original))
		copy(entries, original)
		actual := entries.Reverse()

		assert.Equal(t, actual[0], original[2])
		assert.Equal(t, actual[1], original[1])
		assert.Equal(t, actual[2], original[0])
	})

	t.Run("single", func(t *testing.T) {
		t.Parallel()

		entries := source.LogEntries{
			getFakeLogEntry(),
		}

		assert.Len(t, entries, 1)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		entries := source.LogEntries{}

		assert.Empty(t, entries)
	})
}

func getFakeLogEntry() source.LogEntry {
	return source.LogEntry{
		Time:    "time",
		Level:   source.LevelUnknown,
		Message: "message",
		Line:    []byte(`{"hello":"world"}`),
	}
}

func TestLogEntriesFilter(t *testing.T) {
	t.Parallel()

	term := "special MESSAGE to search by in the test: " + t.Name()

	logEntry := getFakeLogEntry()
	logEntry.Message = term
	logEntry.Line = json.RawMessage(`{"message": "` + term + `"}`)

	logEntries := source.LogEntries{
		getFakeLogEntry(),
		logEntry,
		getFakeLogEntry(),
	}

	t.Run("all", func(t *testing.T) {
		t.Parallel()

		assert.Len(t, logEntries.Filter(""), len(logEntries))
	})

	t.Run("found_exact", func(t *testing.T) {
		t.Parallel()

		filtered := logEntries.Filter(term)
		if assert.Len(t, filtered, 1) {
			assert.Equal(t, logEntry, filtered[0])
		}
	})

	t.Run("found_ignore_case", func(t *testing.T) {
		t.Parallel()

		filtered := logEntries.Filter(strings.ToUpper(term))
		if assert.Len(t, filtered, 1) {
			assert.Equal(t, logEntry, filtered[0])
		}
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()

		filtered := logEntries.Filter(term + " - not found!")
		assert.Empty(t, filtered)
	})
}
