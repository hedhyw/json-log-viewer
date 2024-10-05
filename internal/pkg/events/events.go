package events

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
)

type (
	// LogEntriesUpdateMsg is an event about successfully updated log entries.
	LogEntriesUpdateMsg source.LazyLogEntries
	LogEntriesEOF       struct{}

	// ErrorOccuredMsg is a generic error event.
	ErrorOccuredMsg struct{ Err error }

	// OpenJSONRowRequestedMsg is an event to request extended JSON view
	// for the given row.
	OpenJSONRowRequestedMsg struct {
		// LogEntries include all log entities.
		LogEntries source.LazyLogEntries

		// Index of the row.
		Index int
	}
)

// OpenJSONRowRequested implements tea.Cmd. It creates OpenJSONRowRequestedMsg.
func OpenJSONRowRequested(logEntries source.LazyLogEntries, index int) func() tea.Msg {
	return func() tea.Msg {
		return OpenJSONRowRequestedMsg{
			LogEntries: logEntries,
			Index:      index,
		}
	}
}

// ShowError is an event about occurred error.
func ShowError(err error) func() tea.Msg {
	return func() tea.Msg {
		return ErrorOccuredMsg{Err: err}
	}
}

// HelpKeyClicked is a trigger to display detailed help.
func HelpKeyClicked() tea.Msg {
	return tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'?'},
	}
}

// EscKeyClicked is an "Esc" key event.
func EscKeyClicked() tea.Msg {
	return tea.KeyMsg{Type: tea.KeyEsc}
}

// EnterKeyClicked implements tea.Cmd. It creates a message indicating 'Enter' has been clicked.
func EnterKeyClicked() tea.Msg {
	return tea.KeyMsg{Type: tea.KeyEnter}
}

// ArrowRightKeyClicked implements tea.Cmd. It creates a message indicating 'arrow-right' has been clicked.
func ArrowRightKeyClicked() tea.Msg {
	return tea.KeyMsg{Type: tea.KeyRight}
}

// FilterKeyClicked implements tea.Cmd. It creates a message indicating 'f' has been clicked.
func FilterKeyClicked() tea.Msg {
	return tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'f'},
	}
}
