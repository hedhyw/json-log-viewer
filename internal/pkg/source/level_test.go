package source_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
)

func TestParseLevel(t *testing.T) {
	t.Parallel()

	testCases := [...]struct {
		Input    string
		Expected source.Level
	}{{
		Input:    "",
		Expected: source.LevelUnknown,
	}, {
		Input:    "INFO",
		Expected: source.LevelInfo,
	}, {
		Input:    "debug",
		Expected: source.LevelDebug,
	}, {
		Input:    "info",
		Expected: source.LevelInfo,
	}, {
		Input:    "WRN",
		Expected: source.LevelWarning,
	}, {
		Input:    "erR",
		Expected: source.LevelError,
	}, {
		Input:    "error",
		Expected: source.LevelError,
	}, {
		Input:    "panic",
		Expected: source.LevelPanic,
	}, {
		Input:    "fatal",
		Expected: source.LevelFatal,
	}, {
		Input:    "trace",
		Expected: source.LevelTrace,
	}, {
		Input:    "verbose",
		Expected: source.LevelTrace,
	}, {
		Input:    "  Unknown\t\n",
		Expected: source.Level("unknown"),
	}}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.Input, func(t *testing.T) {
			t.Parallel()

			actual := source.ParseLevel(testCase.Input)
			assert.Equal(t, testCase.Expected, actual)
		})
	}
}
