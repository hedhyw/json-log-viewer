package main

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetStdinSource(t *testing.T) {
	t.Parallel()

	t.Run("ModeNamedPipe", func(t *testing.T) {
		t.Parallel()

		content := t.Name() + "\n"

		file := fakeFile{
			Reader: bytes.NewReader([]byte(content)),
			StatFileInfo: fakeFileInfo{
				FileMode: os.ModeNamedPipe,
			},
		}

		input, err := getStdinReader(file)
		require.NoError(t, err)

		data, err := io.ReadAll(input)
		require.NoError(t, err)
		assert.Equal(t, content, string(data))
	})

	t.Run("ModeCharDevice", func(t *testing.T) {
		t.Parallel()

		file := fakeFile{
			Reader: bytes.NewReader([]byte(t.Name() + "\n")),
			StatFileInfo: fakeFileInfo{
				FileMode: os.ModeCharDevice,
			},
		}

		input, err := getStdinReader(file)
		require.NoError(t, err)

		data, err := io.ReadAll(input)
		require.NoError(t, err)
		assert.Empty(t, data)
	})

	t.Run("Stat_error", func(t *testing.T) {
		t.Parallel()

		// nolint: err113 // Test.
		errStat := errors.New(t.Name())

		file := fakeFile{ErrStat: errStat}

		_, err := getStdinReader(file)
		require.Error(t, err)
		require.ErrorIs(t, err, errStat)
	})
}

type fakeFile struct {
	io.Closer
	io.Reader

	StatFileInfo os.FileInfo
	ErrStat      error
}

// Stat implements fs.File.
func (f fakeFile) Stat() (os.FileInfo, error) {
	return f.StatFileInfo, f.ErrStat
}

type fakeFileInfo struct {
	fs.FileInfo
	FileMode fs.FileMode
}

// Mode implements fs.FileInfo.
func (f fakeFileInfo) Mode() fs.FileMode {
	return f.FileMode
}
