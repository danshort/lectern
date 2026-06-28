// Package config loads lectern's optional per-user TOML configuration.
//
// The file lives at $XDG_CONFIG_HOME/lectern/config.toml, falling back to
// ~/.config/lectern/config.toml when XDG_CONFIG_HOME is unset, on all
// platforms (lectern is a CLI, so terminal users expect ~/.config). The file is
// optional: a missing file yields defaults with no error; a malformed file
// yields defaults plus an error the caller surfaces as a warning. Unrecognized
// keys are ignored so future additions and typos degrade gracefully.
package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// EditorConfig controls how the active artifact is opened (see ResolveOpener).
type EditorConfig struct {
	// OpenWith is one of: "" or "$EDITOR" (use $EDITOR, falling back to vi),
	// "system" (the OS default handler), or any other string (an editor command).
	OpenWith string `toml:"open_with"`
}

// Config is the parsed user configuration. New sections/keys can be added over
// time; the loader tolerates unknown keys.
type Config struct {
	Editor EditorConfig `toml:"editor"`
}

// Path returns the resolved config file path, honoring $XDG_CONFIG_HOME and
// otherwise using ~/.config.
func Path() (string, error) {
	if x := os.Getenv("XDG_CONFIG_HOME"); x != "" {
		return filepath.Join(x, "lectern", "config.toml"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "lectern", "config.toml"), nil
}

// Load reads the config from its resolved path. A missing file (or an
// unresolvable home dir) yields a default Config with a nil error; a malformed
// file yields a default Config and a non-nil error for the caller to warn about.
func Load() (Config, error) {
	path, err := Path()
	if err != nil {
		return Config{}, nil // can't resolve a path → defaults, not fatal
	}
	return LoadFile(path)
}

// LoadFile reads and decodes a config file at an explicit path. It is the
// testable core of Load.
func LoadFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Config{}, nil
		}
		return Config{}, fmt.Errorf("%s: %w", path, err)
	}
	var c Config
	// Unknown keys are reported via the returned MetaData.Undecoded(), which we
	// intentionally ignore — forward compatibility over strictness.
	if _, err := toml.Decode(string(data), &c); err != nil {
		return Config{}, fmt.Errorf("%s: %w", path, err)
	}
	return c, nil
}

// configTemplate seeds a new config file: all settings commented out (so it
// parses to defaults) but documented, so a user editing it sees the options.
const configTemplate = `# lectern configuration — see https://github.com/danshort/lectern#configuration

[editor]
# How pressing "e" opens the active artifact:
#   "$EDITOR" (default) — your $EDITOR, falling back to vi, in the terminal
#   "system"            — the OS default app (open / xdg-open / start), detached
#   "nvim", "code --wait", … — any editor command, in the terminal
# open_with = "$EDITOR"
`

// EnsureFile returns the resolved config path, creating the parent directory
// and a commented starter file when none exists yet. Used by the TUI's
// "edit config" action so the user always lands in a documented file.
func EnsureFile() (string, error) {
	path, err := Path()
	if err != nil {
		return "", err
	}
	switch _, statErr := os.Stat(path); {
	case statErr == nil:
		return path, nil
	case !errors.Is(statErr, fs.ErrNotExist):
		return "", statErr
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(path, []byte(configTemplate), 0o644); err != nil {
		return "", err
	}
	return path, nil
}
