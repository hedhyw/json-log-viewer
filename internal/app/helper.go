package app

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
)

func getColumns(width int) []table.Column {
	const (
		widthTime  = 30
		widthLevel = 10
	)

	return []table.Column{
		{Title: "Time", Width: widthTime},
		{Title: "Level", Width: widthLevel},
		{Title: "Message", Width: width - widthTime - widthLevel},
	}
}

func removeClearSequence(value string) string {
	// https://github.com/charmbracelet/lipgloss/issues/144
	return strings.ReplaceAll(value, "\x1b[0", "\x1b[39")
}
