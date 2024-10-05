package main

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"testing"

	"github.com/hedhyw/json-log-viewer/internal/app"
	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
	"github.com/hedhyw/json-log-viewer/internal/pkg/tests"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunAppVersion(t *testing.T) {
	t.Parallel()

	var outputBuf bytes.Buffer

	err := runApp(applicationArguments{
		Stdout:       &outputBuf,
		PrintVersion: true,
	})
	require.NoError(t, err)

	assert.Contains(t, outputBuf.String(), version)
}

func TestRunAppRunProgramFailed(t *testing.T) {
	t.Parallel()

	fileName := tests.RequireCreateFile(t, []byte(t.Name()))

	err := runApp(applicationArguments{
		Args: []string{fileName},
		RunProgram: func(*tea.Program) (tea.Model, error) {
			return nil, tests.ErrTest
		},
	})
	require.ErrorIs(t, err, tests.ErrTest)
}

func TestRunAppRunProgramReadConfigInvalid(t *testing.T) {
	t.Parallel()

	configPath := tests.RequireCreateFile(t, []byte("invalid config"))

	err := runApp(applicationArguments{
		ConfigPath: configPath,
		RunProgram: func(*tea.Program) (tea.Model, error) {
			t.Fatal("Should not run")

			return app.NewModel("", config.GetDefaultConfig(), version), nil
		},
	})
	require.Error(t, err)
}

func TestRunAppUnexpectedNumberOfArgs(t *testing.T) {
	t.Parallel()

	err := runApp(applicationArguments{
		Args: []string{"1", "2", "3"},
	})
	require.Error(t, err)
}

func TestRunAppReadFileSuccess(t *testing.T) {
	t.Parallel()

	fileName := tests.RequireCreateFile(t, []byte(t.Name()))

	var isStarted bool

	err := runApp(applicationArguments{
		Args: []string{fileName},
		RunProgram: func(p *tea.Program) (tea.Model, error) {
			assert.NotNil(t, p)
			isStarted = true

			return app.NewModel("", config.GetDefaultConfig(), version), nil
		},
	})
	require.NoError(t, err)

	assert.True(t, isStarted)
}

func TestRunAppReadFileNotFound(t *testing.T) {
	t.Parallel()

	err := runApp(applicationArguments{
		Args: []string{t.Name() + "not found"},
		RunProgram: func(*tea.Program) (tea.Model, error) {
			t.Fatal("Should not run")

			return app.NewModel("", config.GetDefaultConfig(), version), nil
		},
	})
	require.Error(t, err)
}

func TestRunAppReadStdinSuccess(t *testing.T) {
	t.Parallel()

	fakeStdin := fakeFile{
		Reader: bytes.NewReader([]byte(t.Name())),
		StatFileInfo: fakeFileInfo{
			FileMode: os.ModeNamedPipe,
		},
	}

	var isStarted bool

	err := runApp(applicationArguments{
		Args:  []string{},
		Stdin: fakeStdin,
		RunProgram: func(p *tea.Program) (tea.Model, error) {
			assert.NotNil(t, p)
			isStarted = true

			return app.NewModel("", config.GetDefaultConfig(), version), nil
		},
	})
	require.NoError(t, err)

	assert.True(t, isStarted)
}

func TestRunAppReadStdinStatFailed(t *testing.T) {
	t.Parallel()

	fakeStdin := fakeFile{
		Reader:  bytes.NewReader([]byte(t.Name())),
		ErrStat: tests.ErrTest,
	}

	err := runApp(applicationArguments{
		Args:  []string{},
		Stdin: fakeStdin,
		RunProgram: func(*tea.Program) (tea.Model, error) {
			t.Fatal("Should not run")

			return app.NewModel("", config.GetDefaultConfig(), version), nil
		},
	})
	require.ErrorIs(t, err, tests.ErrTest)
}

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
