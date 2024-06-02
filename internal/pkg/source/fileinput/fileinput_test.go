package fileinput_test

import (
	"context"
	"io"
	"testing"

	"github.com/hedhyw/json-log-viewer/assets"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source/fileinput"
	"github.com/hedhyw/json-log-viewer/internal/pkg/tests"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileInputString(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	expected := assets.ExampleJSONLog()
	testFile := tests.RequireCreateFile(t, expected)

	t.Run("ReadCloser", func(t *testing.T) {
		t.Parallel()

		input := fileinput.New(testFile)

		readCloser, err := input.ReadCloser(ctx)
		require.NoError(t, err)

		t.Cleanup(func() { assert.NoError(t, readCloser.Close()) })

		actual, err := io.ReadAll(readCloser)
		require.NoError(t, err)

		assert.Equal(t, expected, actual)
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()

		input := fileinput.New(testFile)

		assert.Equal(t, testFile, input.String())
	})

	t.Run("NotFound", func(t *testing.T) {
		t.Parallel()

		input := fileinput.New("not_found_for_" + t.Name())

		_, err := input.ReadCloser(ctx)
		require.Error(t, err)
	})
}
