package source_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source"

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

			assert.Equal(t, "Hello World", fieldKindToValue[config.FieldKindMessage], fieldKindToValue)
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
		testCase := testCase

		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			cfg := config.GetDefaultConfig()

			actual := source.ParseLogEntry(json.RawMessage(testCase.JSON), cfg)

			testCase.Assert(t, getFieldKindToValue(cfg, actual.Fields))
		})
	}
}

func TestLogEntryRow(t *testing.T) {
	t.Parallel()

	entry := getFakeLogEntry()
	row := entry.Row()

	assert.Equal(t, []string(row), entry.Fields)
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
		Fields: []string{
			"time",
			source.LevelUnknown.String(),
			"message",
		},
		Line: []byte(`{"hello":"world"}`),
	}
}

func TestLogEntriesFilter(t *testing.T) {
	t.Parallel()

	term := "special MESSAGE to search by in the test: " + t.Name()

	logEntry := getFakeLogEntry()
	logEntry.Fields = append(logEntry.Fields, term)
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

func getFieldKindToValue(cfg *config.Config, entries []string) map[config.FieldKind]string {
	fieldKindToValue := make(map[config.FieldKind]string, len(entries))

	for i, f := range cfg.Fields {
		fieldKindToValue[f.Kind] = entries[i]
	}

	return fieldKindToValue
}

type TimeFormattingTestCase struct {
	TestName       string
	JSON           string
	ExpectedOutput string
}

func getTimestampFormattingConfig(fieldKind config.FieldKind) *config.Config {
	return &config.Config{
		Path: config.PathDefault,
		Fields: []config.Field{{
			Title:      "Time",
			Kind:       fieldKind,
			References: []string{"$.timestamp", "$.time", "$.t", "$.ts"},
			Width:      30,
		}},
	}
}
func TestSecondTimeFormatting(t *testing.T) {
	t.Parallel()

	cfg := getTimestampFormattingConfig(config.FieldKindSecondTime)

	expectedOutput := time.Unix(1, 0).UTC().Format(time.RFC3339)

	secondsTestCases := [...]TimeFormattingTestCase{{
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
	}}

	for _, testCase := range secondsTestCases {
		testCase := testCase
		t.Run(testCase.TestName, func(t *testing.T) {
			t.Parallel()
			actual := source.ParseLogEntry(json.RawMessage(testCase.JSON), cfg)
			assert.Equal(t, testCase.ExpectedOutput, actual.Fields[0])
		})
	}
}

func TestMillisecondTimeFormatting(t *testing.T) {
	t.Parallel()

	cfg := getTimestampFormattingConfig(config.FieldKindMilliTime)

	expectedOutput := time.Unix(2, 0).UTC().Format(time.RFC3339)

	millisecondTestCases := [...]TimeFormattingTestCase{{
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
		testCase := testCase
		t.Run(testCase.TestName, func(t *testing.T) {
			t.Parallel()
			actual := source.ParseLogEntry(json.RawMessage(testCase.JSON), cfg)
			assert.Equal(t, testCase.ExpectedOutput, actual.Fields[0])
		})
	}
}

func TestMicrosecondTimeFormatting(t *testing.T) {
	t.Parallel()

	cfg := getTimestampFormattingConfig(config.FieldKindMicroTime)

	expectedOutput := time.Unix(4, 0).UTC().Format(time.RFC3339)

	microsecondTestCases := [...]TimeFormattingTestCase{{
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
		testCase := testCase
		t.Run(testCase.TestName, func(t *testing.T) {
			t.Parallel()
			actual := source.ParseLogEntry(json.RawMessage(testCase.JSON), cfg)
			assert.Equal(t, testCase.ExpectedOutput, actual.Fields[0])
		})
	}
}

func TestNumericKindTimeFormatting(t *testing.T) {
	t.Parallel()

	cfg := getTimestampFormattingConfig(config.FieldKindNumericTime)

	numericKindCases := [...]TimeFormattingTestCase{{
		TestName:       "Date passthru",
		JSON:           `{"timestamp":"2023-10-08 20:00:00"}`,
		ExpectedOutput: "2023-10-08 20:00:00",
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
	}}

	for _, testCase := range numericKindCases {
		testCase := testCase
		t.Run(testCase.TestName, func(t *testing.T) {
			t.Parallel()
			actual := source.ParseLogEntry(json.RawMessage(testCase.JSON), cfg)
			assert.Equal(t, testCase.ExpectedOutput, actual.Fields[0])
		})
	}
}
