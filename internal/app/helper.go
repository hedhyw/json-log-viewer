package app

import "github.com/charmbracelet/bubbles/table"

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
