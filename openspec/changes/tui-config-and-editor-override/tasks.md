## 1. Config loader

- [x] 1.1 Add the `github.com/BurntSushi/toml` dependency (`go get`; tidy)
- [x] 1.2 Add a user-config loader (e.g. `internal/config`): resolve the path (`$XDG_CONFIG_HOME/lectern/config.toml` else `~/.config/lectern/config.toml`), decode TOML into a struct with an `Editor.OpenWith string`, ignore unknown keys
- [x] 1.3 Missing file → zero-value/default config, no error; malformed file → return the parse error so the caller can warn
- [x] 1.4 Tests: XDG path vs `~/.config` resolution, missing file → defaults, malformed → error surfaced, unknown keys ignored, `editor.open_with` parsed

## 2. Opener resolution

- [x] 2.1 Add an opener resolver mapping `open_with` → (mode, argv): unset/`"$EDITOR"` → terminal (`$EDITOR` fields else `vi`); `"system"` → detached (`open`/`xdg-open`/`start` by GOOS); other → terminal (command fields)
- [x] 2.2 Tests for the resolver across the three modes and the platform handler selection

## 3. Wire into the TUI

- [x] 3.1 In `cmd/lectern/main.go`, load the user config; on a malformed-config error print a warning to stderr and continue with defaults; thread the resolved editor preference into the model (`ui.New` / `ui.NewSinglePath`)
- [x] 3.2 In `internal/ui/viewer.go` `openInEditor`, branch on mode: terminal → `tea.ExecProcess` (as today); detached → verify the handler exists, `Start()` it without yielding the terminal, then trigger the normal reload
- [x] 3.3 Surface launch failures (opener not found / exec error) via `m.errMsg` instead of always returning `editorReturnMsg{}` (`viewer.go:170`)
- [x] 3.4 Apply the same resolution to the spec-view editor open path (`ModeViewingSpec`)
- [x] 3.5 Add a `c` keybinding (index + change-viewer) that opens the user config in the editor, creating a documented starter file if missing (`config.EnsureFile`); advertise it in the help overlay; test `EnsureFile`

## 4. Verification

- [x] 4.1 `make test` / `go test -race ./...` passes (config loader, resolver, and existing editor-launch tests still green)
- [x] 4.2 `go vet ./...` clean; `gofmt` clean; `make lint` if available
- [ ] 4.3 Manual: no config → `e` still opens `$EDITOR` (terminal) as before
- [ ] 4.4 Manual: `open_with = "system"` → `e` opens the artifact in the default GUI app without freezing the TUI; saving reflects on the next reload
- [ ] 4.5 Manual: `open_with = "nvim"` → `e` opens nvim in-terminal and resumes on exit
- [ ] 4.6 Manual: malformed `config.toml` → warning on stderr, TUI still launches with defaults
- [x] 4.7 Update README/DEVELOPING with the config file location and the `[editor] open_with` setting
