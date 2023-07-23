package app

// IsErrorShown indicates that err is shown on the screen.
func (m Model) IsErrorShown() bool {
	return m.err != nil
}

// IsJSONShown indicates that extended JSON view is shown on the screen.
func (m Model) IsJSONShown() bool {
	return m.jsonView != nil
}

// IsJSONShown indicates that the main list is shown on the screen.
func (m Model) IsTableShown() bool {
	return !m.IsJSONShown()
}

// IsFilterShown indicates that the filter is shown on the screen.
func (m Model) IsFilterShown() bool {
	return m.textInputShown
}

// IsFiltered indicates that the results are filtered.
func (m Model) IsFiltered() bool {
	return len(m.allLogEntries) != len(m.filteredLogEntries)
}
