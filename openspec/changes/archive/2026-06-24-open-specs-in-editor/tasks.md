## 1. Extract shared editor-launch helper

- [x] 1.1 Add an `openInEditor(path string) tea.Cmd` method on `*Model` (in `internal/ui/viewer.go` or a new `internal/ui/editor.go`) containing the `$EDITOR`-split / `vi`-fallback / `exec.Command` / `tea.ExecProcess` logic and the existing `#nosec G204 G702` annotation.
- [x] 1.2 Replace the inline `case "e":` body in `internal/ui/viewer.go` to call `m.openInEditor(m.artifactPath())`, preserving the `tabAvailable` / non-empty-path guards.

## 2. Resolve and open the viewed spec

- [x] 2.1 Add a `currentSpecPath() string` method on `*Model` returning `filepath.Join(m.root, "openspec", "specs", m.projectSpecs[m.specViewer.Cursor].Name, "spec.md")`, or `""` when `m.specViewer.Cursor` is out of range.
- [x] 2.2 Add `case "e":` to `updateSpec` in `internal/ui/spec.go`: resolve `currentSpecPath()` and, if non-empty, return `m, m.openInEditor(path)` (no full/focus branching needed).

## 3. Reload spec content on editor return

- [x] 3.1 In the `editorReturnMsg` handler in `internal/ui/update.go`, when `m.mode == ModeViewingSpec`, reload `m.projectSpecs` via `m.loader.LoadProjectSpecsFrom(m.root)` and clamp `m.specViewer.Cursor` if the list shrank, before `m.loadViewport()`.
- [x] 3.2 Confirm the existing change-reload branch remains gated to change modes so spec edits don't trigger a change reload.

## 4. Help bar

- [x] 4.1 In `renderHelpBar()` (`internal/ui/view.go`), add `e: edit` to the `ModeViewingSpec` full-spec help text.
- [x] 4.2 Add `e: edit` to the `ModeViewingSpec` requirement-focus help text.

## 5. Tests

- [x] 5.1 Add a test for `currentSpecPath()` covering a valid cursor (correct path) and an out-of-range cursor (empty string).
- [x] 5.2 Add a test asserting `updateSpec` returns a non-nil command for `e` when a spec is viewed, and that `$EDITOR` unset falls back to `vi` (table-driven, mirroring existing editor-launch tests).
- [x] 5.3 Add a test that an `editorReturnMsg` in `ModeViewingSpec` refreshes `m.projectSpecs` content from disk and stays in `ModeViewingSpec`.
- [x] 5.4 Run `gofmt`, `go vet`, and the full test suite; ensure no regressions in existing editor-launch tests.

## 6. Validate change

- [x] 6.1 Run `openspec validate open-specs-in-editor` and resolve any issues.
