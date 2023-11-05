package app

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"

	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
)

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
