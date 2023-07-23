package tests

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// RequireCreateFile is a helper that create a temporary file and deletes
// it at the end of the test.
func RequireCreateFile(tb testing.TB, content []byte) string {
	tb.Helper()

	f, err := os.CreateTemp("", "json_log_viewer_test")
	require.NoError(tb, err)

	defer func() { assert.NoError(tb, f.Close()) }()

	_, err = f.Write(content)
	require.NoError(tb, err)

	name := f.Name()
	tb.Cleanup(func() { assert.NoError(tb, os.Remove(name)) })

	return name
}
