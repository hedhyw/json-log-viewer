package source

import "strings"

// Level of the logs entity.
type Level string

// String implement fmt.Stringer interface.
func (l Level) String() string {
	return strings.ToLower(string(l))
}

// Possible log levels.
const (
	LevelUnknown Level = "none"
)
