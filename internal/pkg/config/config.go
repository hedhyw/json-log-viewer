package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/hedhyw/jsoncjson"
)

// PathDefault is a fake path to the default config.
const PathDefault = "default"

// Config contains application customization settings.
type Config struct {
	// Path to the config.
	Path string `json:"-"`

	Fields []Field `json:"fields" validate:"min=1"`
}

// FieldKind describes the type of the log field.
type FieldKind string

// Possible kinds.
const (
	FieldKindTime        FieldKind = "time"
	FieldKindNumericTime FieldKind = "numerictme"
	FieldKindSecondTime  FieldKind = "secondtme"
	FieldKindMilliTime   FieldKind = "millitme"
	FieldKindMicroTime   FieldKind = "microtme"
	FieldKindMessage     FieldKind = "message"
	FieldKindLevel       FieldKind = "level"
	FieldKindAny         FieldKind = "any"
)

// Field customization.
type Field struct {
	Title      string    `json:"title" validate:"required,min=1,max=32"`
	Kind       FieldKind `json:"kind" validate:"required,oneof=time message numerictme secondtme millitme microtme level any"`
	References []string  `json:"ref" validate:"min=1,dive,required"`
	Width      int       `json:"width" validate:"min=0"`
}

// GetDefaultConfig returns the configuration with default values.
func GetDefaultConfig() *Config {
	return &Config{
		Path: "default",
		Fields: []Field{{
			Title:      "Time",
			Kind:       FieldKindTime,
			References: []string{"$.timestamp", "$.time", "$.t", "$.ts"},
			Width:      30,
		}, {
			Title:      "Level",
			Kind:       FieldKindLevel,
			References: []string{"$.level", "$.lvl", "$.l"},
			Width:      10,
		}, {
			Title:      "Message",
			Kind:       FieldKindMessage,
			References: []string{"$.message", "$.msg", "$.error", "$.err"},
		}},
	}
}

// Read config from the given paths. From higher priority to lower priority.
func Read(paths ...string) (*Config, error) {
	cfg, err := readConfigFromPaths(paths...)
	if err != nil {
		return nil, fmt.Errorf("reading from paths: %w", err)
	}

	err = validator.New().Struct(cfg)
	if err != nil {
		return nil, fmt.Errorf("validating config: %s: %w", cfg.Path, err)
	}

	return cfg, nil
}

func readConfigFromPaths(paths ...string) (*Config, error) {
	for _, p := range paths {
		_, err := os.Stat(p)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}

			return nil, fmt.Errorf("checking config: %w", err)
		}

		cfg, err := readConfigFromFile(p)
		if err != nil {
			return nil, fmt.Errorf("reading config from file: %w", err)
		}

		return cfg, nil
	}

	return GetDefaultConfig(), nil
}

func readConfigFromFile(path string) (cfg *Config, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("os opening: %w", err)
	}

	defer func() { err = errors.Join(err, file.Close()) }()

	err = json.NewDecoder(
		jsoncjson.NewReader(file),
	).Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("decoding json: %w", err)
	}

	cfg.Path = path

	return cfg, nil
}
