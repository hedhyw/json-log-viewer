package events

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
)

type (
	// LogEntriesLoadedMsg is an event about successfully loaded log entries.
	LogEntriesLoadedMsg source.LogEntries

	// ErrorOccuredMsg is a generic error event.
	ErrorOccuredMsg struct{ Err error }

	// OpenJSONRowRequestedMsg is an event to request extended JSON view
	// for the given row.
	OpenJSONRowRequestedMsg struct {
		// LogEntries include all log entities.
		LogEntries source.LogEntries

		// Index of the row.
		Index int
	}

	// ViewRowsReloadRequestedMsg is an event to start reloading of logs.
	ViewRowsReloadRequestedMsg struct{}
)

// OpenJSONRowRequested implements tea.Cmd. It creates OpenJSONRowRequestedMsg.
func OpenJSONRowRequested(logEntries source.LogEntries, index int) func() tea.Msg {
	return func() tea.Msg {
		return OpenJSONRowRequestedMsg{
			LogEntries: logEntries,
			Index:      index,
		}
	}
}

// ViewRowsReloadRequested implements tea.Cmd. It creates ViewRowsReloadRequestedMsg.
func ViewRowsReloadRequested() tea.Msg {
	return ViewRowsReloadRequestedMsg{}
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

// BackKeyClicked implements tea.Cmd. It creates a message indicating 'Esc' has been clicked.
func BackKeyClicked() tea.Msg {
	return tea.KeyMsg{Type: tea.KeyEscape}
}
