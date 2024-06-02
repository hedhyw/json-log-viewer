package app_test

import (
	"errors"
	"testing"

	"github.com/charmbracelet/bubbles/cursor"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/require"

	"github.com/hedhyw/json-log-viewer/internal/app"
	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source/fileinput"
	"github.com/hedhyw/json-log-viewer/internal/pkg/tests"
)

const testVersion = "v0.0.1"

func newTestModel(tb testing.TB, content []byte) tea.Model {
	tb.Helper()

	testFile := tests.RequireCreateFile(tb, content)

	model := app.NewModel(fileinput.New(testFile), config.GetDefaultConfig(), testVersion)
	model = handleUpdate(model, model.Init()())

	return model
}

func handleUpdate(model tea.Model, msg tea.Msg) tea.Model {
	model, cmd := model.Update(msg)

	const limit = 10
	var i int

	var cmdsBatch []tea.Cmd

	if cmd != nil {
		cmdsBatch = append(cmdsBatch, cmd)
	}

	for len(cmdsBatch) > 0 && i < limit {
		i++

		cmd = cmdsBatch[len(cmdsBatch)-1]
		cmdsBatch = cmdsBatch[:len(cmdsBatch)-1]

		if msg = cmd(); msg == nil {
			return model
		}

		if _, ok := msg.(cursor.BlinkMsg); ok {
			break
		}

		if batch, ok := msg.(tea.BatchMsg); ok {
			cmdsBatch = append(cmdsBatch, batch...)

			continue
		}

		if model, cmd = model.Update(msg); cmd != nil {
			cmdsBatch = append(cmdsBatch, cmd)
		}
	}

	return model
}

func requireCmdMsg(tb testing.TB, expected tea.Msg, cmd tea.Cmd) {
	tb.Helper()

	require.NotNil(tb, cmd)

	msg := cmd()

	if batch, ok := msg.(tea.BatchMsg); ok {
		for _, cmd := range batch {
			msg := cmd()

			tb.Logf("%T: %v\n", msg, msg)

			if msg == expected {
				return
			}
		}

		require.Failf(tb, "batch message doesn't include expected msg", "%+v", batch)
	} else {
		require.Equal(tb, expected, msg)
	}
}

func getTestError() error {
	// nolint: goerr113 // It is a test.
	return errors.New("error description")
}
