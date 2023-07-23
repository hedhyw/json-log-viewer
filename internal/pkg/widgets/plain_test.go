package widgets_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"

	"github.com/hedhyw/json-log-viewer/internal/pkg/widgets"
)

const text = "hello world"

func TestPlainLogModelInit(t *testing.T) {
	model, _ := widgets.NewPlainLogModel(text, getFakeTeaWindowSizeMsg())

	cmd := model.Init()
	assert.Nil(t, cmd)
}

func TestPlainLogModelUpdateTeaWindowSizeMsg(t *testing.T) {
	windowSize := getFakeTeaWindowSizeMsg()
	model, _ := widgets.NewPlainLogModel(text, windowSize)

	windowSize.Height++
	windowSize.Width++

	actual, _ := model.Update(windowSize)
	if assert.NotNil(t, actual) {
		assert.NotEqual(t, actual, model)
	}
}

func TestPlainLogModelView(t *testing.T) {
	model, _ := widgets.NewPlainLogModel(text, getFakeTeaWindowSizeMsg())

	actual := model.View()
	assert.Contains(t, actual, text)
}

func getFakeTeaWindowSizeMsg() tea.WindowSizeMsg {
	return tea.WindowSizeMsg{Width: 100, Height: 100}
}
