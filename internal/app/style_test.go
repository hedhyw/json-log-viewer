package app

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"

	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
)

func TestGetColorForLogLevel(t *testing.T) {
	t.Parallel()

	testCases := [...]struct {
		Level    source.Level
		Expected lipgloss.Color
	}{{
		Level:    source.Level(""),
		Expected: "",
	}, {
		Level:    source.LevelUnknown,
		Expected: "",
	}, {
		Level:    source.Level("custom"),
		Expected: "",
	}, {
		Level:    source.LevelTrace,
		Expected: colorMagenta,
	}, {
		Level:    source.LevelDebug,
		Expected: colorYellow,
	}, {
		Level:    source.LevelInfo,
		Expected: colorGreen,
	}, {
		Level:    source.LevelWarning,
		Expected: colorOrange,
	}, {
		Level:    source.LevelError,
		Expected: colorRed,
	}, {
		Level:    source.LevelFatal,
		Expected: colorRed,
	}, {
		Level:    source.LevelPanic,
		Expected: colorRed,
	}}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.Level.String(), func(t *testing.T) {
			t.Parallel()

			actual := getColorForLogLevel(testCase.Level)
			assert.Equal(t, testCase.Expected, actual)
		})
	}
}
