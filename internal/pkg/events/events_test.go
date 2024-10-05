package events_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"

	"github.com/hedhyw/json-log-viewer/internal/pkg/events"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
	"github.com/hedhyw/json-log-viewer/internal/pkg/tests"
)

func TestEvents(t *testing.T) {
	t.Parallel()

	testCases := [...]struct {
		Name     string
		Actual   tea.Msg
		Expected tea.Msg
	}{{
		Name:   "OpenJSONRowRequested",
		Actual: events.OpenJSONRowRequested(source.LazyLogEntries{}, 0)(),
		Expected: events.OpenJSONRowRequestedMsg{
			LogEntries: source.LazyLogEntries{},
			Index:      0,
		},
	}, {
		Name:     "ShowError",
		Actual:   events.ShowError(tests.ErrTest)(),
		Expected: events.ErrorOccuredMsg{Err: tests.ErrTest},
	}, {
		Name:     "HelpKeyClicked",
		Actual:   events.HelpKeyClicked(),
		Expected: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}},
	}, {
		Name:     "EscKeyClicked",
		Actual:   events.EscKeyClicked(),
		Expected: tea.KeyMsg{Type: tea.KeyEsc},
	}, {
		Name:     "EnterKeyClicked",
		Actual:   events.EnterKeyClicked(),
		Expected: tea.KeyMsg{Type: tea.KeyEnter},
	}, {
		Name:     "ArrowRightKeyClicked",
		Actual:   events.ArrowRightKeyClicked(),
		Expected: tea.KeyMsg{Type: tea.KeyRight},
	}, {
		Name:     "FilterKeyClicked",
		Actual:   events.FilterKeyClicked(),
		Expected: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}},
	}}

	for _, testCase := range testCases {
		assert.Equal(t,
			testCase.Expected,
			testCase.Actual,
			testCase.Name,
		)
	}
}
