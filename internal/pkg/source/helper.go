package source

import (
	"unicode"
)

func normalizeJSON(input []byte) []byte {
	out := make([]byte, 0, len(input))

	for _, r := range string(input) {
		if unicode.IsPrint(r) {
			out = append(out, []byte(string(r))...)
		}
	}

	return out
}
