package source_test

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/stretchr/testify/require"

	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
	"github.com/hedhyw/json-log-viewer/internal/pkg/tests"

	"github.com/stretchr/testify/assert"
)

func TestParseLogEntryDefault(t *testing.T) {
	t.Parallel()

	testCases := [...]struct {
		Name   string
		JSON   string
		Assert func(tb testing.TB, fieldKindToValue map[config.FieldKind]string)
	}{{
		Name: "plain_log",
		JSON: "Hello World",
		Assert: func(tb testing.TB, fieldKindToValue map[config.FieldKind]string) {
			tb.Helper()

			assert.Equal(t, "Hello World\n", fieldKindToValue[config.FieldKindMessage], fieldKindToValue)
			assert.Equal(t, "-", fieldKindToValue[config.FieldKindLevel], fieldKindToValue)
			assert.Equal(t, "-", fieldKindToValue[config.FieldKindNumericTime], fieldKindToValue)
		},
	}, {
		Name: "time_number",
		JSON: `{"time":1}`,
		Assert: func(tb testing.TB, fieldKindToValue map[config.FieldKind]string) {
			tb.Helper()

			assert.Equal(t, time.Unix(1, 0).UTC().Format(time.RFC3339), fieldKindToValue[config.FieldKindNumericTime], fieldKindToValue)
		},
	}, {
		Name: "timestamp_number",
		JSON: `{"timestamp":1}`,
		Assert: func(tb testing.TB, fieldKindToValue map[config.FieldKind]string) {
			tb.Helper()

			assert.Equal(t, time.Unix(1, 0).UTC().Format(time.RFC3339), fieldKindToValue[config.FieldKindNumericTime], fieldKindToValue)
		},
	}, {
		Name: "ts_number",
		JSON: `{"ts":1}`,
		Assert: func(tb testing.TB, fieldKindToValue map[config.FieldKind]string) {
			tb.Helper()

			assert.Equal(t, time.Unix(1, 0).UTC().Format(time.RFC3339), fieldKindToValue[config.FieldKindNumericTime], fieldKindToValue)
		},
	}, {
		Name: "ts_int_seconds_as_string",
		JSON: `{"ts":"1"}`,
		Assert: func(tb testing.TB, fieldKindToValue map[config.FieldKind]string) {
			tb.Helper()

			assert.Equal(t, time.Unix(1, 0).UTC().Format(time.RFC3339), fieldKindToValue[config.FieldKindNumericTime], fieldKindToValue)
		},
	}, {
		Name: "ts_float_seconds_as_string",
		JSON: `{"ts":"1.29333384"}`,
		Assert: func(tb testing.TB, fieldKindToValue map[config.FieldKind]string) {
			tb.Helper()

			assert.Equal(t, time.Unix(1, 0).UTC().Format(time.RFC3339), fieldKindToValue[config.FieldKindNumericTime], fieldKindToValue)
		},
	}, {
		Name: "ts_",
		JSON: `{"ts":"1.29333384"}`,
		Assert: func(tb testing.TB, fieldKindToValue map[config.FieldKind]string) {
			tb.Helper()

			assert.Equal(t, time.Unix(1, 0).UTC().Format(time.RFC3339), fieldKindToValue[config.FieldKindNumericTime], fieldKindToValue)
		},
	}, {
		Name: "time_text",
		JSON: `{"time":"1970-01-01T00:00:00.00"}`,
		Assert: func(tb testing.TB, fieldKindToValue map[config.FieldKind]string) {
			tb.Helper()

			assert.Equal(t,
				"1970-01-01T00:00:00.00",
				fieldKindToValue[config.FieldKindNumericTime],
				fieldKindToValue,
			)
		},
	}, {
		Name: "timestamp_text",
		JSON: `{"timestamp":"1970-01-01T00:00:00.00"}`,
		Assert: func(tb testing.TB, fieldKindToValue map[config.FieldKind]string) {
			tb.Helper()

			assert.Equal(t,
				"1970-01-01T00:00:00.00",
				fieldKindToValue[config.FieldKindNumericTime],
				fieldKindToValue,
			)
		},
	}, {
		Name: "ts_text",
		JSON: `{"ts":"1970-01-01T00:00:00.00"}`,
		Assert: func(tb testing.TB, fieldKindToValue map[config.FieldKind]string) {
			tb.Helper()

			assert.Equal(t,
				"1970-01-01T00:00:00.00",
				fieldKindToValue[config.FieldKindNumericTime],
				fieldKindToValue,
			)
		},
	}, {
		Name: "message",
		JSON: `{"message":"message"}`,
		Assert: func(tb testing.TB, fieldKindToValue map[config.FieldKind]string) {
			tb.Helper()

			assert.Equal(t,
				"message",
				fieldKindToValue[config.FieldKindMessage],
				fieldKindToValue,
			)
		},
	}, {
		Name: "msg",
		JSON: `{"msg":"msg"}`,
		Assert: func(tb testing.TB, fieldKindToValue map[config.FieldKind]string) {
			tb.Helper()

			assert.Equal(t,
				"msg",
				fieldKindToValue[config.FieldKindMessage],
				fieldKindToValue,
			)
		},
	}, {
		Name: "message_special_rune",
		JSON: `{"message":"mes` + string(rune(1)) + `sage"}`,
		Assert: func(tb testing.TB, fieldKindToValue map[config.FieldKind]string) {
			tb.Helper()

			assert.Equal(t,
				"message",
				fieldKindToValue[config.FieldKindMessage],
				fieldKindToValue,
			)
		},
	}, {
		Name: "error",
		JSON: `{"error":"error"}`,
		Assert: func(tb testing.TB, fieldKindToValue map[config.FieldKind]string) {
			tb.Helper()

			assert.Equal(t,
				"error",
				fieldKindToValue[config.FieldKindMessage],
				fieldKindToValue,
			)
		},
	}, {
		Name: "err",
		JSON: `{"err":"err"}`,
		Assert: func(tb testing.TB, fieldKindToValue map[config.FieldKind]string) {
			tb.Helper()

			assert.Equal(t,
				"err",
				fieldKindToValue[config.FieldKindMessage],
				fieldKindToValue,
			)
		},
	}, {
		Name: "special",
		JSON: `{"msg":"\u0008"}`,
		Assert: func(tb testing.TB, fieldKindToValue map[config.FieldKind]string) {
			tb.Helper()

			assert.Empty(t,
				fieldKindToValue[config.FieldKindMessage],
				fieldKindToValue,
			)
		},
	}, {
		Name: "level",
		JSON: `{"level":"INFO"}`,
		Assert: func(tb testing.TB, fieldKindToValue map[config.FieldKind]string) {
			tb.Helper()

			assert.Equal(t,
				"info",
				fieldKindToValue[config.FieldKindLevel],
				fieldKindToValue,
			)
		},
	}}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			cfg := config.GetDefaultConfig()

			actual := parseTableRow(t, testCase.JSON, cfg)
			testCase.Assert(t, getFieldKindToValue(cfg, actual))
		})
	}
}

func TestLogEntryRow(t *testing.T) {
	t.Parallel()

	entry := getFakeLogEntry()
	row := entry.Row()

	assert.Equal(t, []string(row), entry.Fields)
}

func getFakeLogEntry() source.LogEntry {
	return source.LogEntry{
		Fields: []string{
			"time",
			source.LevelUnknown.String(),
			"message",
		},
		Line: []byte(`{"hello":"world"}`),
	}
}

func TestLazyLogEntriesFilter(t *testing.T) {
	t.Parallel()

	term := "special MESSAGE to search by in the test: " + t.Name()

	logs := fmt.Sprintf(`
{"hello":"world"}
{"message", "%s"}
{"hello":"world"}
`, term)

	createEntries := func(tb testing.TB) (source.LazyLogEntries, source.LazyLogEntry) {
		source, err := source.Reader(bytes.NewReader([]byte(logs)), config.GetDefaultConfig())
		require.NoError(t, err)

		tb.Cleanup(func() { assert.NoError(tb, source.Close()) })

		logEntries, err := source.ParseLogEntries()
		require.NoError(t, err)

		logEntry := logEntries.Entries[1]

		return logEntries, logEntry
	}

	t.Run("all", func(t *testing.T) {
		t.Parallel()

		logEntries, _ := createEntries(t)

		assert.Len(t, logEntries.Entries, logEntries.Len())
	})

	t.Run("found_exact", func(t *testing.T) {
		t.Parallel()

		logEntries, logEntry := createEntries(t)

		filtered, err := logEntries.Filter(term)
		require.NoError(t, err)

		if assert.Len(t, filtered.Entries, 1) {
			assert.Equal(t, logEntry, filtered.Entries[0])
		}
	})

	t.Run("found_ignore_case", func(t *testing.T) {
		t.Parallel()

		logEntries, logEntry := createEntries(t)

		filtered, err := logEntries.Filter(strings.ToUpper(term))
		require.NoError(t, err)

		if assert.Len(t, filtered.Entries, 1) {
			assert.Equal(t, logEntry, filtered.Entries[0])
		}
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		logEntries, _ := createEntries(t)

		filtered, err := logEntries.Filter("")
		require.NoError(t, err)
		assert.Len(t, filtered.Entries, logEntries.Len())
	})

	t.Run("not_found", func(t *testing.T) {
		t.Parallel()

		logEntries, _ := createEntries(t)

		filtered, err := logEntries.Filter(term + " - not found!")
		require.NoError(t, err)

		assert.Empty(t, filtered.Entries)
	})

	t.Run("seeker_failed", func(t *testing.T) {
		t.Parallel()

		logEntries, _ := createEntries(t)

		fileName := tests.RequireCreateFile(t, []byte(""))

		f, err := os.Open(fileName)
		require.NoError(t, err)
		require.NoError(t, f.Close())

		logEntries.Seeker = f

		_, err = logEntries.Filter(term + " - not found!")
		require.Error(t, err)
	})
}

func TestSecondTimeFormatting(t *testing.T) {
	t.Parallel()

	expectedOutput := time.Unix(1, 0).UTC().Format(time.RFC3339)

	secondsTestCases := [...]timeFormattingTestCase{{
		TestName:       "Seconds (float)",
		JSON:           `{"timestamp":1.0}`,
		ExpectedOutput: expectedOutput,
	}, {
		TestName:       "Seconds (int)",
		JSON:           `{"timestamp":1}`,
		ExpectedOutput: expectedOutput,
	}, {
		TestName:       "Seconds (float as string)",
		JSON:           `{"timestamp":"1.0"}`,
		ExpectedOutput: expectedOutput,
	}, {
		TestName:       "Seconds (int as string)",
		JSON:           `{"timestamp":"1"}`,
		ExpectedOutput: expectedOutput,
	}, {
		TestName:       "Seconds (int as string)",
		JSON:           `{"timestamp":"x"}`,
		ExpectedOutput: `x`,
	}}

	for _, testCase := range secondsTestCases {
		t.Run(testCase.TestName, func(t *testing.T) {
			t.Parallel()

			cfg := getTimestampFormattingConfig(config.FieldKindSecondTime, testCase.Format)

			actual := parseTableRow(t, testCase.JSON, cfg)
			assert.Equal(t, testCase.ExpectedOutput, actual[0])
		})
	}
}

func TestMillisecondTimeFormatting(t *testing.T) {
	t.Parallel()

	expectedOutput := time.Unix(2, 0).UTC().Format(time.RFC3339)

	millisecondTestCases := [...]timeFormattingTestCase{{
		TestName:       "Milliseconds (float)",
		JSON:           `{"timestamp":2000.0}`,
		ExpectedOutput: expectedOutput,
	}, {
		TestName:       "Milliseconds (int)",
		JSON:           `{"timestamp":2000}`,
		ExpectedOutput: expectedOutput,
	}, {
		TestName:       "Milliseconds (float as string)",
		JSON:           `{"timestamp":"2000.0"}`,
		ExpectedOutput: expectedOutput,
	}, {
		TestName:       "Milliseconds (int as string)",
		JSON:           `{"timestamp":"2000"}`,
		ExpectedOutput: expectedOutput,
	}}

	for _, testCase := range millisecondTestCases {
		t.Run(testCase.TestName, func(t *testing.T) {
			t.Parallel()

			cfg := getTimestampFormattingConfig(config.FieldKindMilliTime, testCase.Format)

			actual := parseTableRow(t, testCase.JSON, cfg)
			assert.Equal(t, testCase.ExpectedOutput, actual[0])
		})
	}
}

func TestMicrosecondTimeFormatting(t *testing.T) {
	t.Parallel()

	expectedOutput := time.Unix(4, 0).UTC().Format(time.RFC3339)

	microsecondTestCases := [...]timeFormattingTestCase{{
		TestName:       "Microseconds (float)",
		JSON:           `{"timestamp":4000000.0}`,
		ExpectedOutput: expectedOutput,
	}, {
		TestName:       "Microseconds (int)",
		JSON:           `{"timestamp":4000000}`,
		ExpectedOutput: expectedOutput,
	}, {
		TestName:       "Microseconds (float as string)",
		JSON:           `{"timestamp":"4000000.0"}`,
		ExpectedOutput: expectedOutput,
	}, {
		TestName:       "Microseconds (int as string)",
		JSON:           `{"timestamp":"4000000"}`,
		ExpectedOutput: expectedOutput,
	}}

	for _, testCase := range microsecondTestCases {
		t.Run(testCase.TestName, func(t *testing.T) {
			t.Parallel()

			cfg := getTimestampFormattingConfig(config.FieldKindMicroTime, testCase.Format)

			actual := parseTableRow(t, testCase.JSON, cfg)
			assert.Equal(t, testCase.ExpectedOutput, actual[0])
		})
	}
}

func TestFormattingUnknown(t *testing.T) {
	t.Parallel()

	cfg := getTimestampFormattingConfig(config.FieldKind("unknown"), config.DefaultTimeFormat)

	actual := parseTableRow(t, `{"timestamp": 1}`, cfg)
	assert.Equal(t, "1", actual[0])
}

func TestFormattingAny(t *testing.T) {
	t.Parallel()

	cfg := getTimestampFormattingConfig(config.FieldKindAny, config.DefaultTimeFormat)

	actual := parseTableRow(t, `{"timestamp": 1}`, cfg)
	assert.Equal(t, "1", actual[0])
}

func TestNumericKindTimeFormatting(t *testing.T) {
	t.Parallel()

	numericKindCases := [...]timeFormattingTestCase{{
		TestName:       "Date passthru",
		JSON:           `{"timestamp":"2023-10-08 20:00:00"}`,
		ExpectedOutput: "2023-10-08T20:00:00Z",
	}, {
		TestName:       "RFC1123 passthru",
		JSON:           `{"@timestamp":"Mon, 02 Jan 2006 15:04:05 MST"}`,
		ExpectedOutput: "2006-01-02T15:04:05Z",
	}, {
		TestName:       "Non-date string",
		JSON:           `{"timestamp":"-"}`,
		ExpectedOutput: "-",
	}, {
		TestName:       "Seconds as int",
		JSON:           `{"timestamp":4000000}`,
		ExpectedOutput: time.Unix(4000000, 0).UTC().Format(time.RFC3339),
	}, {
		TestName:       "Seconds as int string",
		JSON:           `{"timestamp":"4000000"}`,
		ExpectedOutput: time.Unix(4000000, 0).UTC().Format(time.RFC3339),
	}, {
		TestName:       "Seconds as float",
		JSON:           `{"timestamp":4000000.1}`,
		ExpectedOutput: time.Unix(4000000, 0).UTC().Format(time.RFC3339),
	}, {
		TestName:       "Seconds as float string",
		JSON:           `{"timestamp":"4000000.1"}`,
		ExpectedOutput: time.Unix(4000000, 0).UTC().Format(time.RFC3339),
	}, {
		TestName:       "11 character int is in milliseconds",
		JSON:           `{"timestamp":12345678900}`,
		ExpectedOutput: time.Unix(12345678, 0).UTC().Format(time.RFC3339),
	}, {
		TestName:       "float with 11 digits before the decimal is milliseconds",
		JSON:           `{"timestamp":12345678000000.222}`,
		ExpectedOutput: time.Unix(12345678, 0).UTC().Format(time.RFC3339),
	}, {
		TestName:       "14 character int is in microseconds",
		JSON:           `{"timestamp":12345678900000}`,
		ExpectedOutput: time.Unix(12345678, 0).UTC().Format(time.RFC3339),
	}, {
		TestName:       "float with 14 digits before the decimal is microseconds",
		JSON:           `{"timestamp":12345678900000.222}`,
		ExpectedOutput: time.Unix(12345678, 0).UTC().Format(time.RFC3339),
	}, {
		TestName:       "max_int64",
		JSON:           fmt.Sprintf(`{"timestamp":"%d"}`, math.MaxInt64),
		ExpectedOutput: strconv.Itoa(math.MaxInt64),
	}, {
		TestName:       "negative",
		JSON:           `{"timestamp":"-1"}`,
		ExpectedOutput: "-1",
	}}

	for _, testCase := range numericKindCases {
		t.Run(testCase.TestName, func(t *testing.T) {
			t.Parallel()

			cfg := getTimestampFormattingConfig(config.FieldKindNumericTime, testCase.Format)

			actual := parseTableRow(t, testCase.JSON, cfg)
			assert.Equal(t, testCase.ExpectedOutput, actual[0])
		})
	}
}

func TestLazyLogEntryLength(t *testing.T) {
	t.Parallel()

	entry := t.Name() + "\n"

	logEntry := parseLazyLogEntry(t, entry, config.GetDefaultConfig())
	assert.Equal(t, len(entry), logEntry.Length())
}

func TestLazyLogEntryLine(t *testing.T) {
	t.Parallel()

	entry := t.Name() + "\n"

	logEntry := parseLazyLogEntry(t, entry, config.GetDefaultConfig())
	assert.Equal(t, len(entry), logEntry.Length())

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		fileName := tests.RequireCreateFile(t, []byte(entry))

		f, err := os.Open(fileName)
		require.NoError(t, err)

		t.Cleanup(func() { assert.NoError(t, f.Close()) })

		actual, err := logEntry.Line(f)
		require.NoError(t, err)

		assert.Equal(t, entry, string(actual))
	})

	t.Run("failed", func(t *testing.T) {
		t.Parallel()

		fileName := tests.RequireCreateFile(t, []byte(entry))

		f, err := os.Open(fileName)
		require.NoError(t, err)
		require.NoError(t, f.Close())

		_, err = logEntry.Line(f)
		require.Error(t, err)
	})
}

func TestLazyLogEntryLogEntry(t *testing.T) {
	t.Parallel()

	entry := t.Name() + "\n"
	cfg := config.GetDefaultConfig()

	logEntry := parseLazyLogEntry(t, entry, config.GetDefaultConfig())
	assert.Equal(t, len(entry), logEntry.Length())

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		fileName := tests.RequireCreateFile(t, []byte(entry))

		f, err := os.Open(fileName)
		require.NoError(t, err)

		t.Cleanup(func() { assert.NoError(t, f.Close()) })

		actual := logEntry.LogEntry(f, cfg)
		require.NoError(t, actual.Error)
		assert.Equal(t, entry, string(actual.Line))
	})

	t.Run("failed", func(t *testing.T) {
		t.Parallel()

		fileName := tests.RequireCreateFile(t, []byte(entry))

		f, err := os.Open(fileName)
		require.NoError(t, err)
		require.NoError(t, f.Close())

		actual := logEntry.LogEntry(f, cfg)
		require.Error(t, actual.Error)
	})
}

func TestTimeFormat(t *testing.T) {
	t.Parallel()

	logDate := time.Date(
		2000, // Year.
		time.January,
		2, // Day.
		3, // Hour.
		4, // Minutes.
		5, // Seconds.
		0, // Nanoseconds.
		time.UTC,
	)

	jsonContent := fmt.Sprintf(`{"timestamp":"%d"}`, logDate.Unix())

	numericKindCases := [...]timeFormattingTestCase{{
		TestName:       "RFC3339",
		JSON:           jsonContent,
		ExpectedOutput: logDate.Format(time.RFC3339),
		Format:         time.RFC3339,
	}, {
		TestName:       "RFC1123",
		JSON:           jsonContent,
		ExpectedOutput: logDate.Format(time.RFC1123),
		Format:         time.RFC1123,
	}, {
		TestName:       "TimeOnly",
		JSON:           jsonContent,
		ExpectedOutput: logDate.Format(time.TimeOnly),
		Format:         time.TimeOnly,
	}, {
		TestName:       "TimeOnly",
		JSON:           jsonContent,
		ExpectedOutput: "invalid",
		Format:         "invalid",
	}}

	for _, testCase := range numericKindCases {
		t.Run(testCase.TestName, func(t *testing.T) {
			t.Parallel()

			cfg := getTimestampFormattingConfig(config.FieldKindSecondTime, testCase.Format)

			actual := parseTableRow(t, testCase.JSON, cfg)
			assert.Equal(t, testCase.ExpectedOutput, actual[0])
		})
	}
}

func parseLazyLogEntry(tb testing.TB, value string, cfg *config.Config) source.LazyLogEntry {
	tb.Helper()

	source, err := source.Reader(strings.NewReader(value), cfg)
	require.NoError(tb, err)

	tb.Cleanup(func() { assert.NoError(tb, source.Close()) })

	logEntries, err := source.ParseLogEntries()
	require.NoError(tb, err)
	require.Equal(tb, 1, logEntries.Len())

	return logEntries.Entries[0]
}

func parseTableRow(tb testing.TB, value string, cfg *config.Config) table.Row {
	tb.Helper()

	source, err := source.Reader(strings.NewReader(value+"\n"), cfg)
	require.NoError(tb, err)

	tb.Cleanup(func() { assert.NoError(tb, source.Close()) })

	logEntries, err := source.ParseLogEntries()
	require.NoError(tb, err)
	require.Equal(tb, 1, logEntries.Len(), value)

	return logEntries.Row(cfg, 0)
}

func getFieldKindToValue(cfg *config.Config, entries []string) map[config.FieldKind]string {
	fieldKindToValue := make(map[config.FieldKind]string, len(entries))

	for i, f := range cfg.Fields {
		fieldKindToValue[f.Kind] = entries[i]
	}

	return fieldKindToValue
}

type timeFormattingTestCase struct {
	TestName       string
	JSON           string
	ExpectedOutput string
	Format         string
}

func getTimestampFormattingConfig(fieldKind config.FieldKind, format string) *config.Config {
	cfg := config.GetDefaultConfig()

	var timeFormat *string

	if format != "" {
		timeFormat = &format
	}

	cfg.Fields = []config.Field{{
		Title:      "Time",
		Kind:       fieldKind,
		References: []string{"$.timestamp", "$.time", "$.t", "$.ts", "$[\"@timestamp\"]"},
		Width:      30,
		TimeFormat: timeFormat,
	}}

	return cfg
}
