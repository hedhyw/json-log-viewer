package app

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"

	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
)

// Possible colors.
const (
	colorMagenta lipgloss.Color = "13"
	colorYellow  lipgloss.Color = "11"
	colorGreen   lipgloss.Color = "10"
	colorOrange  lipgloss.Color = "214"
	colorRed     lipgloss.Color = "9"
)

func getTableStyles() table.Styles {
	tableStyles := table.DefaultStyles()
	tableStyles.Header = tableStyles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	tableStyles.Selected = tableStyles.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	return tableStyles
}

func getBaseStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))
}

func getFooterStyle() lipgloss.Style {
	return lipgloss.NewStyle().Height(footerSize).PaddingLeft(2)
}

func (m Model) getLogLevelStyle(baseStyle lipgloss.Style, rowID int) lipgloss.Style {
	if rowID < 0 || rowID >= len(m.filteredLogEntries) {
		return baseStyle
	}

	color := getColorForLogLevel(m.filteredLogEntries[rowID].Level)
	if color == "" {
		return baseStyle
	}

	return baseStyle.Copy().Foreground(color)
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
