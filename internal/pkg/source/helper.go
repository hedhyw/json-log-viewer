package source

import (
	"strconv"
	"strings"

	"github.com/valyala/fastjson"
)

func extractTime(value *fastjson.Value) string {
	timeValue := extractValue(value, "timestamp", "time", "t", "ts")
	if timeValue != "" {
		return formatMessage(strings.TrimSpace(timeValue))
	}

	return "-"
}

func extractLevel(value *fastjson.Value) Level {
	level := extractValue(value, "level", "lvl")

	return ParseLevel(formatMessage(level))
}

func extractValue(value *fastjson.Value, keys ...string) string {
	for _, k := range keys {
		element := value.Get(k)

		text := string(element.GetStringBytes())
		if text != "" {
			return text
		}

		number := element.GetInt()
		if number != 0 {
			return strconv.Itoa(number)
		}
	}

	return ""
}

func extractMessage(value *fastjson.Value) string {
	message := extractValue(value, "message", "msg", "error", "err")
	if message != "" {
		return strings.TrimSpace(message)
	}

	return strings.TrimSpace(value.String())
}
