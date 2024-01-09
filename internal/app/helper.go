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

type helper struct {
	Application
}

func (h helper) LoadEntries() tea.Msg {
	logEntries, err := source.LoadLogsFromFile(
		h.Path,
		h.Config,
	)
	if err != nil {
		return events.ErrorOccuredMsg{Err: err}
	}

	return events.LogEntriesLoadedMsg(logEntries)
}

func (h helper) getLogLevelStyle(
	logEntries source.LogEntries,
	baseStyle lipgloss.Style,
	rowID int,
) lipgloss.Style {
	if rowID < 0 || rowID >= len(logEntries) {
		return baseStyle
	}

	color := getColorForLogLevel(h.getLogLevelFromLogEntry(logEntries[rowID]))
	if color == "" {
		return baseStyle
	}

	return baseStyle.Copy().Foreground(color)
}

// Update application state.
func (h helper) Update(msg tea.Msg) helper {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		h.LastWindowSize = msg
	}

	return h
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

func (h helper) getLogLevelFromLogEntry(logEntry source.LogEntry) source.Level {
	return source.Level(getFieldByKind(h.Config, config.FieldKindLevel, logEntry))
}

func (h helper) handleErrorOccuredMsg(msg events.ErrorOccuredMsg) (tea.Model, tea.Cmd) {
	return initializeModel(newStateError(h.Application, msg.Err))
}

func (h helper) handleLogEntriesLoadedMsg(
	msg events.LogEntriesLoadedMsg,
) (tea.Model, tea.Cmd) {
	return initializeModel(newStateViewLogs(
		h.Application,
		source.LogEntries(msg)),
	)
}

func (h helper) handleOpenJSONRowRequestedMsg(
	msg events.OpenJSONRowRequestedMsg,
	previousState state,
) (tea.Model, tea.Cmd) {
	if msg.Index < 0 || msg.Index >= len(msg.LogEntries) {
		return previousState, nil
	}

	logEntry := msg.LogEntries[msg.Index]

	return initializeModel(newStateViewRow(
		h.Application,
		logEntry,
		previousState,
	))
}

func (h helper) handleKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch {
	case h.isQuitKeyMap(msg):
		return tea.Quit
	case key.Matches(msg, defaultKeys.Filter):
		return events.FilterKeyClicked
	case h.isEnterKeyMap(msg):
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
