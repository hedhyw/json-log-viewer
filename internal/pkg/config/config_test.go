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

	"github.com/docker/go-units"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
	jsonEncoder.SetIndent("", "  ")

	if err := jsonEncoder.Encode(&cfg); err != nil {
		log.Fatal(err)
	}

	fmt.Println(buf.String())
	// Output:
	// {
	//   "fields": [
	//     {
	//       "title": "Time",
	//       "kind": "numerictime",
	//       "ref": [
	//         "$.timestamp",
	//         "$.time",
	//         "$.t",
	//         "$.ts"
	//       ],
	//       "width": 30,
	//       "time_format": "2006-01-02T15:04:05Z07:00"
	//     },
	//     {
	//       "title": "Level",
	//       "kind": "level",
	//       "ref": [
	//         "$.level",
	//         "$.lvl",
	//         "$.l"
	//       ],
	//       "width": 10
	//     },
	//     {
	//       "title": "Message",
	//       "kind": "message",
	//       "ref": [
	//         "$.message",
	//         "$.msg",
	//         "$.error",
	//         "$.err"
	//       ],
	//       "width": 0
	//     }
	//   ],
	//   "customLevelMapping": {
	//     "10": "trace",
	//     "20": "debug",
	//     "30": "info",
	//     "40": "warn",
	//     "50": "error",
	//     "60": "fatal"
	//   },
	//   "maxFileSizeBytes": 2000000000
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

func TestByteSize(t *testing.T) {
	t.Parallel()

	testCases := [...]struct {
		Value    string
		Expected config.ByteSize
	}{{
		Value:    `"1k"`,
		Expected: units.KB,
	}, {
		Value:    `"1m"`,
		Expected: units.MB,
	}, {
		Value:    `"1.5m"`,
		Expected: units.MB * 1.5,
	}, {
		Value:    `"1g"`,
		Expected: units.GB,
	}, {
		Value:    `"1t"`,
		Expected: units.TB,
	}, {
		Value:    `"1p"`,
		Expected: units.PB,
	}, {
		Value:    "1",
		Expected: 1,
	}, {
		Value:    "12345",
		Expected: 12345,
	}}

	for _, testCase := range testCases {
		var actual config.ByteSize

		err := json.Unmarshal([]byte(testCase.Value), &actual)
		if assert.NoError(t, err, testCase.Value) {
			assert.Equal(t, testCase.Expected, actual, testCase.Value)
		}
	}
}

func TestByteSizeParseFailed(t *testing.T) {
	t.Parallel()

	t.Run("invalid_number", func(t *testing.T) {
		t.Parallel()

		var value config.ByteSize

		err := json.Unmarshal([]byte(`"123.123.123"`), &value)
		require.Error(t, err)
	})

	t.Run("invalid_suffix", func(t *testing.T) {
		t.Parallel()

		var value config.ByteSize

		err := json.Unmarshal([]byte(`"123X"`), &value)
		require.Error(t, err)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		var value config.ByteSize

		err := json.Unmarshal([]byte(`""`), &value)
		require.Error(t, err)
	})
}
