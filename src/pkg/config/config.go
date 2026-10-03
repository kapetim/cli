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

// Config is the effective configuration.
type Config struct {
	Markdown Markdown `json:"markdown"`
	Linters  []string `json:"linters"`
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
