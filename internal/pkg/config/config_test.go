package config_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"

	"github.com/hedhyw/json-log-viewer/internal/pkg/config"
	"github.com/hedhyw/json-log-viewer/internal/pkg/tests"
)

func TestReadDefault(t *testing.T) {
	t.Parallel()

	cfg, err := config.Read()
	if assert.NoError(t, err) {
		assert.Equal(t, config.PathDefault, cfg.Path)

		def := config.GetDefaultConfig()
		assert.ElementsMatch(t, cfg.Fields, def.Fields)
	}
}

func TestReadNotFound(t *testing.T) {
	t.Parallel()

	cfg, err := config.Read("not_found_" + strconv.FormatInt(time.Now().UnixNano(), 10))
	if assert.NoError(t, err) {
		assert.Equal(t, config.PathDefault, cfg.Path)
	}
}

func TestReadPriority(t *testing.T) {
	t.Parallel()

	configJSON := tests.RequireEncodeJSON(t, config.GetDefaultConfig())
	fileFirst := tests.RequireCreateFile(t, configJSON)
	fileSecond := tests.RequireCreateFile(t, configJSON)

	cfg, err := config.Read(fileFirst, fileSecond)
	if assert.NoError(t, err) {
		assert.Equal(t, fileFirst, cfg.Path)
	}
}

func TestReadValidated(t *testing.T) {
	t.Parallel()

	cfg := config.GetDefaultConfig()
	cfg.Fields = nil

	configJSON := tests.RequireEncodeJSON(t, cfg)
	configFile := tests.RequireCreateFile(t, configJSON)

	_, err := config.Read(configFile)
	if assert.Error(t, err) {
		assert.ErrorAs(t, err, &validator.ValidationErrors{})
	}
}

func TestReadInvalidJSON(t *testing.T) {
	t.Parallel()

	configFile := tests.RequireCreateFile(t, []byte("-"))

	_, err := config.Read(configFile)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, io.ErrUnexpectedEOF)
	}
}

func TestReadDirectory(t *testing.T) {
	t.Parallel()

	_, err := config.Read(".")
	assert.Error(t, err)
}

func ExampleGetDefaultConfig() {
	cfg := config.GetDefaultConfig()

	var buf bytes.Buffer

	jsonEncoder := json.NewEncoder(&buf)
	jsonEncoder.SetIndent("", "\t")

	if err := jsonEncoder.Encode(&cfg); err != nil {
		log.Fatal(err)
	}

	fmt.Println(buf.String())
	// Output:
	// {
	// 	"fields": [
	// 		{
	// 			"title": "Time",
	// 			"kind": "numerictime",
	// 			"ref": [
	// 				"$.timestamp",
	// 				"$.time",
	// 				"$.t",
	// 				"$.ts"
	// 			],
	// 			"width": 30
	// 		},
	// 		{
	// 			"title": "Level",
	// 			"kind": "level",
	// 			"ref": [
	// 				"$.level",
	// 				"$.lvl",
	// 				"$.l"
	// 			],
	// 			"width": 10
	// 		},
	// 		{
	// 			"title": "Message",
	// 			"kind": "message",
	// 			"ref": [
	// 				"$.message",
	// 				"$.msg",
	// 				"$.error",
	// 				"$.err"
	// 			],
	// 			"width": 0
	// 		}
	// 	]
	// }
}

func TestValidateField(t *testing.T) {
	t.Parallel()

	testCases := [...]struct {
		Name    string
		Apply   func(value *config.Field)
		IsValid bool
	}{{
		Name:    "ok",
		Apply:   func(*config.Field) {},
		IsValid: true,
	}, {
		Name: "unset_title",
		Apply: func(value *config.Field) {
			value.Title = ""
		},
		IsValid: false,
	}, {
		Name: "short_title",
		Apply: func(value *config.Field) {
			value.Title = "."
		},
		IsValid: true,
	}, {
		Name: "almost_long_title",
		Apply: func(value *config.Field) {
			value.Title = strings.Repeat(".", 32)
		},
		IsValid: true,
	}, {
		Name: "long_title",
		Apply: func(value *config.Field) {
			value.Title = strings.Repeat(".", 33)
		},
		IsValid: false,
	}, {
		Name: "unset_references",
		Apply: func(value *config.Field) {
			value.References = []string{}
		},
		IsValid: false,
	}, {
		Name: "empty_reference",
		Apply: func(value *config.Field) {
			value.References = []string{""}
		},
		IsValid: false,
	}, {
		Name: "kind_any",
		Apply: func(value *config.Field) {
			value.Kind = config.FieldKindAny
		},
		IsValid: true,
	}, {
		Name: "kind_level",
		Apply: func(value *config.Field) {
			value.Kind = config.FieldKindLevel
		},
		IsValid: true,
	}, {
		Name: "kind_message",
		Apply: func(value *config.Field) {
			value.Kind = config.FieldKindMessage
		},
		IsValid: true,
	}, {
		Name: "kind_time",
		Apply: func(value *config.Field) {
			value.Kind = config.FieldKindTime
		},
		IsValid: true,
	}, {
		Name: "kind_numeric_time",
		Apply: func(value *config.Field) {
			value.Kind = config.FieldKindNumericTime
		},
		IsValid: true,
	}, {
		Name: "kind_second_time",
		Apply: func(value *config.Field) {
			value.Kind = config.FieldKindSecondTime
		},
		IsValid: true,
	}, {
		Name: "kind_milli_time",
		Apply: func(value *config.Field) {
			value.Kind = config.FieldKindMilliTime
		},
		IsValid: true,
	}, {
		Name: "kind_micro_time",
		Apply: func(value *config.Field) {
			value.Kind = config.FieldKindMicroTime
		},
		IsValid: true,
	}, {
		Name: "unset_kind",
		Apply: func(value *config.Field) {
			value.Kind = ""
		},
		IsValid: false,
	}, {
		Name: "invalid_kind",
		Apply: func(value *config.Field) {
			value.Kind = "invalid"
		},
		IsValid: false,
	}, {
		Name: "unset_width",
		Apply: func(value *config.Field) {
			value.Width = 0
		},
		IsValid: true,
	}, {
		Name: "small_width",
		Apply: func(value *config.Field) {
			value.Width = 1
		},
		IsValid: true,
	}, {
		Name: "negative_width",
		Apply: func(value *config.Field) {
			value.Width = -1
		},
		IsValid: false,
	}}

	validator := validator.New()

	for _, testCaseNotInParallel := range testCases {
		testCase := testCaseNotInParallel

		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			value := config.Field{
				Title:      "Title",
				Kind:       config.FieldKindAny,
				References: []string{"$.test"},
				Width:      0,
			}

			testCase.Apply(&value)

			err := validator.Struct(value)
			if testCase.IsValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
