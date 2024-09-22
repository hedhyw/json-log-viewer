package app_test

import (
	"strconv"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestLazyTableModelManyRows(t *testing.T) {
	t.Parallel()

	const (
		prefix     = "n"
		countLines = 2_000
	)

	var content []byte

	for i := range countLines {
		content = append(content, prefix+strconv.Itoa(i)...)
		content = append(content, '\n')
	}

	model := newTestModel(t, content)

	step := strings.Count(model.View(), prefix)

	for i := countLines - 1; i > countLines-step; i-- {
		expectedMessage := prefix + strconv.Itoa(i)
		assert.Contains(t, model.View(), expectedMessage)
		model = handleUpdate(model, tea.KeyMsg{Type: tea.KeyDown})
	}

	for i := countLines - step - 1; i >= 0; i-- {
		model = handleUpdate(model, tea.KeyMsg{Type: tea.KeyDown})
		expectedMessage := prefix + strconv.Itoa(i)
		assert.Contains(t, model.View(), expectedMessage)
	}
}

func TestLazyTableModelReverse(t *testing.T) {
	t.Parallel()

	const (
		start          = "START"
		end            = "END"
		keywordReverse = "reverse"
	)

	var (
		keyReverse    = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
		keyGoToBottom = tea.KeyMsg{Type: tea.KeyEnd}
		keyGoToUp     = tea.KeyMsg{Type: tea.KeyHome}

		middle  = strings.Repeat("-\n", 100)
		content = []byte(start + "\n" + middle + end + "\n")
	)

	t.Run("forward", func(t *testing.T) {
		t.Parallel()

		model := newTestModel(t, content)

		model = handleUpdate(model, keyReverse)

		view := model.View()

		// "END" should be at the first half of the screen.
		assert.NotContains(t, view, keywordReverse)
		assert.True(t, strings.Index(view, end) > (len(view)/2), view)
	})

	t.Run("reverse", func(t *testing.T) {
		t.Parallel()

		model := newTestModel(t, content)

		model = handleUpdate(model, keyReverse)
		model = handleUpdate(model, keyReverse)

		view := model.View()

		// "END" should be at the second half of the screen.
		assert.True(t, strings.Index(view, end) < (len(view)/2), view)
		assert.Contains(t, view, keywordReverse)
	})

	t.Run("reverse_default", func(t *testing.T) {
		t.Parallel()

		model := newTestModel(t, content)
		view := model.View()
		assert.Contains(t, view, keywordReverse)
	})

	t.Run("reverse_go_to_bottom", func(t *testing.T) {
		t.Parallel()

		model := newTestModel(t, content)
		view := model.View()
		assert.Contains(t, view, keywordReverse)

		model = handleUpdate(model, keyGoToBottom)

		view = model.View()

		// "START" should be at the second half of the screen.
		assert.True(t, strings.Index(view, start) > (len(view)/2), view)
	})

	t.Run("reverse_go_to_bottom_and_up", func(t *testing.T) {
		t.Parallel()

		model := newTestModel(t, content)
		view := model.View()
		assert.Contains(t, view, keywordReverse)

		model = handleUpdate(model, keyGoToBottom)
		model = handleUpdate(model, keyGoToUp)

		view = model.View()

		assert.Contains(t, view, end)
	})

	t.Run("forwards_go_to_up_and_bottom", func(t *testing.T) {
		t.Parallel()

		model := newTestModel(t, content)
		view := model.View()
		assert.Contains(t, view, keywordReverse)

		model = handleUpdate(model, keyReverse)
		model = handleUpdate(model, keyGoToUp)
		model = handleUpdate(model, keyGoToBottom)

		view = model.View()

		assert.Contains(t, view, end)
	})

	t.Run("forward_go_to_bottom", func(t *testing.T) {
		t.Parallel()

		model := newTestModel(t, content)
		view := model.View()
		assert.Contains(t, view, keywordReverse)

		model = handleUpdate(model, keyGoToUp)

		view = model.View()

		// "START" should be at the second half of the screen.
		assert.True(t, strings.Index(view, start) < (len(view)/2), view)
	})
}
