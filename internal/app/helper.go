package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
	"github.com/hedhyw/json-log-viewer/internal/pkg/events"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
)

func (app *Application) getLogLevelStyle(
	logEntries source.LazyLogEntries,
	baseStyle lipgloss.Style,
	rowID int,
) lipgloss.Style {
	if rowID < 0 || rowID >= logEntries.Len() {
		return baseStyle
	}

	entry := logEntries.Entries[rowID].LogEntry(logEntries.Seeker, app.Config)

	color := getColorForLogLevel(app.getLogLevelFromLogEntry(entry))
	if color == "" {
		return baseStyle
	}

	return baseStyle.Foreground(color)
}

// Update application state.
func (app *Application) Update(msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		app.LastWindowSize = msg
	case events.LogEntriesUpdateMsg:
		app.Entries = source.LazyLogEntries(msg)
	}
}

func getColorForLogLevel(level source.Level) lipgloss.Color {
	switch level {
	case source.LevelTrace:
		return colorMagenta
	case source.LevelDebug:
		return colorYellow
	case source.LevelInfo:
		return colorGreen
	case source.LevelWarning:
		return colorOrange
	case source.LevelError,
		source.LevelFatal,
		source.LevelPanic:
		return colorRed
	default:
		return ""
	}
}

func (app *Application) getLogLevelFromLogEntry(logEntry source.LogEntry) source.Level {
	return source.Level(getFieldByKind(app.Config, config.FieldKindLevel, logEntry))
}

func (app *Application) handleErrorOccuredMsg(msg events.ErrorOccuredMsg) (tea.Model, tea.Cmd) {
	return initializeModel(newStateError(app, msg.Err))
}

func (app *Application) handleInitialLogEntriesLoadedMsg(
	msg events.LogEntriesUpdateMsg,
) (tea.Model, tea.Cmd) {
	return initializeModel(newStateViewLogs(
		app,
		source.LazyLogEntries(msg),
	))
}

func (app *Application) handleOpenJSONRowRequestedMsg(
	msg events.OpenJSONRowRequestedMsg,
	previousState stateModel,
) (tea.Model, tea.Cmd) {
	if msg.Index < 0 || msg.Index >= msg.LogEntries.Len() {
		return previousState, nil
	}

	logEntry := msg.LogEntries.Entries[msg.Index]

	return initializeModel(newStateViewRow(
		logEntry.LogEntry(msg.LogEntries.Seeker, app.Config),
		previousState,
	))
}

func (app *Application) handleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, defaultKeys.Exit):
		return tea.Quit
	case key.Matches(msg, defaultKeys.Filter):
		return events.FilterKeyClicked
	case key.Matches(msg, defaultKeys.Open):
		return events.EnterKeyClicked
	case key.Matches(msg, defaultKeys.ToggleViewArrow):
		return events.ArrowRightKeyClicked
	default:
		return nil
	}
}

func getColumns(width int, cfg *config.Config) []table.Column {
	const minWidth = 10

	flexSpace := width
	flexColumns := 0

	for _, f := range cfg.Fields {
		flexSpace -= f.Width

		if f.Width == 0 {
			flexColumns++
		}
	}

	flexWidth := 0

	if flexColumns != 0 {
		flexWidth = max(minWidth, flexSpace/flexColumns)
	}

	colums := make([]table.Column, 0, len(cfg.Fields))

	for _, f := range cfg.Fields {
		if f.Width == 0 {
			f.Width = flexWidth
		}

		colums = append(colums, table.Column{
			Title: f.Title,
			Width: f.Width,
		})
	}

	return colums
}

func removeClearSequence(value string) string {
	// https://github.com/charmbracelet/lipgloss/issues/144
	return strings.ReplaceAll(value, "\x1b[0", "\x1b[39")
}

func getFieldByKind(
	cfg *config.Config,
	kind config.FieldKind,
	logEntry source.LogEntry,
) string {
	for i, f := range cfg.Fields {
		if f.Kind != kind {
			continue
		}

		if i >= len(logEntry.Fields) {
			return "-"
		}

		return logEntry.Fields[i]
	}

	return ""
}

func batched[T any](m T, cmd tea.Cmd) func(batch []tea.Cmd) (T, []tea.Cmd) {
	return func(batch []tea.Cmd) (T, []tea.Cmd) {
		if cmd != nil {
			batch = append(batch, cmd)
		}

		return m, batch
	}
}

func appendCmd(batch []tea.Cmd, cmd tea.Cmd) []tea.Cmd {
	if cmd != nil {
		batch = append(batch, cmd)
	}

	return batch
}

func initializeModel[T tea.Model](m T) (T, tea.Cmd) {
	return m, m.Init()
}

func modelValue(model tea.Model) string {
	return fmt.Sprintf("%T\n%s", model, model.View())
}
