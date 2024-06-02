package app

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// Component sizes.
const (
	footerSize        = 1
	footerPaddingLeft = 2
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
	return lipgloss.NewStyle().Height(footerSize).PaddingLeft(footerPaddingLeft)
}
