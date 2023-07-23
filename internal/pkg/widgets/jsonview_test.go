package widgets_test

import (
	"testing"

	"github.com/hedhyw/json-log-viewer/internal/pkg/widgets"

	"github.com/stretchr/testify/assert"
)

func TestNewJSONViewModel(t *testing.T) {
	t.Parallel()

	t.Run("plain_text", func(t *testing.T) {
		t.Parallel()

		model, _ := widgets.NewJSONViewModel([]byte(text), getFakeTeaWindowSizeMsg())

		_, ok := model.(widgets.PlainLogModel)
		assert.Truef(t, ok, "actual type: %T", model)
	})

	t.Run("json", func(t *testing.T) {
		t.Parallel()

		model, _ := widgets.NewJSONViewModel(
			[]byte(`{"hello":"world"}`),
			getFakeTeaWindowSizeMsg(),
		)

		_, ok := model.(widgets.PlainLogModel)
		assert.Falsef(t, ok, "actual type: %T", model)
	})
}
