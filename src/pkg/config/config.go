// Package config resolves the effective configuration: embedded defaults
// overlaid with an optional per-repo override.
package config

import (
	_ "embed"
	"encoding/json"
	"os"
	"path/filepath"
)

//go:embed defaults.json
var defaultJSON []byte

// Markdown holds markdown-lint settings.
type Markdown struct {
	MaxCell int `json:"maxCell"`
}

// Filenames holds filename-convention settings for `validate filenames`.
type Filenames struct {
	Readme      string   `json:"readme"`      // upper | lower | any
	Dockerfiles bool     `json:"dockerfiles"` // allow Dockerfile / *.Dockerfile
	Allow       []string `json:"allow"`       // exact uppercase names that are allowed
	Extensions  []string `json:"extensions"`  // allowed file extensions (lowercase)
	Skip        []string `json:"skip"`        // directory names to skip
}

// Config is the effective configuration.
type Config struct {
	Markdown  Markdown  `json:"markdown"`
	Linters   []string  `json:"linters"`
	Filenames Filenames `json:"filenames"`
}

// Defaults returns the embedded defaults.
func Defaults() Config {
	var c Config
	_ = json.Unmarshal(defaultJSON, &c)
	return c
}

// Load overlays an optional `.cli.json` in dir on top of the defaults.
func Load(dir string) (Config, error) {
	c := Defaults()
	b, err := os.ReadFile(filepath.Join(dir, ".cli.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return c, nil
		}
		return c, err
	}
	if err := json.Unmarshal(b, &c); err != nil {
		return Defaults(), err
	}
	return c, nil
}
