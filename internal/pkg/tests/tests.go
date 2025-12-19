package tests

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/hedhyw/semerr/pkg/v1/semerr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ErrTest is a fake constant error to use in tests.
const ErrTest semerr.Error = "test error"

// RequireCreateFile is a helper that create a temporary file and deletes
// it at the end of the test.
func RequireCreateFile(tb testing.TB, content []byte) string {
	tb.Helper()

	f, err := os.CreateTemp(tb.TempDir(), "json_log_viewer_test")
	require.NoError(tb, err)

	defer func() { assert.NoError(tb, f.Close()) }()

	_, err = f.Write(content)
	require.NoError(tb, err)

	name := f.Name()
	tb.Cleanup(func() {
		if _, err := os.Stat(name); err == nil {
			assert.NoError(tb, os.Remove(name))
		}
	})

	return name
}

// RequireEncodeJSON marshals value to JSON.
func RequireEncodeJSON(tb testing.TB, value any) []byte {
	tb.Helper()

	content, err := json.Marshal(value)
	require.NoError(tb, err)

	return content
}

// Context returns a test context with timeout.
func Context(t *testing.T) context.Context {
	t.Helper()

	const defaultTimeout = time.Minute

	deadline, ok := t.Deadline()
	if !ok {
		deadline = time.Now().Add(defaultTimeout)
	}

	ctx, cancel := context.WithDeadline(t.Context(), deadline)
	t.Cleanup(cancel)

	return ctx
}
