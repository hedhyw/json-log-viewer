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

	// EnterKeyClickedMsg is a keyboard event after pressing <Enter>.
	EnterKeyClickedMsg struct{}

	// EnterKeyClickedMsg is a keyboard event after pressing <Right>.
	ArrowRightKeyClickedMsg struct{}

	// FilterKeyClickedMsg is a keyboard event for "Filtering" key.
	FilterKeyClickedMsg struct{}
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

// EnterKeyClicked implements tea.Cmd. It creates EnterKeyClickedMsg.
func EnterKeyClicked() tea.Msg {
	return EnterKeyClickedMsg{}
}

// ArrowRightKeyClicked implements tea.Cmd. It creates ArrowRightKeyClickedMsg.
func ArrowRightKeyClicked() tea.Msg {
	return ArrowRightKeyClickedMsg{}
}

// FilterKeyClicked implements tea.Cmd. It creates FilterKeyClickedMsg.
func FilterKeyClicked() tea.Msg {
	return FilterKeyClickedMsg{}
}

// BackKeyClicked implements tea.Cmd. It triggers a back key click event.
func BackKeyClicked() tea.Msg {
	return tea.KeyMsg{Type: tea.KeyEscape}
}
