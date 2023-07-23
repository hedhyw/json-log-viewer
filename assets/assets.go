package assets

import (
	_ "embed"
)

//go:embed example.log
var exampleJSONLog []byte

// ExampleJSONLog returns a copy of the file "example.log".
func ExampleJSONLog() []byte {
	target := make([]byte, len(exampleJSONLog))

	copy(target, exampleJSONLog)

	return target
}
