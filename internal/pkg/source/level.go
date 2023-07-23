package source

import "strings"

// Level of the logs entity.
type Level string

// ParseLevel parses level from the text value.
func ParseLevel(value string) Level {
	value = strings.ToLower(value)
	value = strings.TrimSpace(value)

	switch {
	case value == "":
		return LevelUnknown
	case strings.HasPrefix(value, "t"),
		strings.HasPrefix(value, "v"): // Verbose.
		return LevelTrace
	case strings.HasPrefix(value, "d"):
		return LevelDebug
	case strings.HasPrefix(value, "e"):
		return LevelError
	case strings.HasPrefix(value, "i"):
		return LevelInfo
	case strings.HasPrefix(value, "w"):
		return LevelWarning
	case strings.HasPrefix(value, "f"):
		return LevelFatal
	case strings.HasPrefix(value, "p"):
		return LevelPanic
	default:
		return Level(value)
	}
}

// String implement fmt.Stringer interface.
func (l Level) String() string {
	return strings.ToLower(string(l))
}

// Possible log levels.
const (
	LevelUnknown Level = "none"
	LevelTrace   Level = "trace"
	LevelDebug   Level = "debug"
	LevelInfo    Level = "info"
	LevelWarning Level = "warn"
	LevelError   Level = "error"
	LevelPanic   Level = "panic"
	LevelFatal   Level = "fatal"
)
