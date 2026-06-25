## 1. Baseline

- [x] 1.1 Run `go test ./...` and confirm the suite is green before any change (records the behavior contract)
- [x] 1.2 Capture the full list of `m.mode = …` assignment sites with `grep -rn 'm\.mode =' internal/ui --include='*.go'` (excluding tests) as the migration checklist

## 2. Sub-struct grouping (viewer + spec; render/layout stay flat)

- [x] 2.1 Define `viewerState` (`changeIdx`, `tab`, `specIdx`) in `model.go`
- [x] 2.2 Replace the flat `changeIdx`/`tab`/`specIdx` fields on `Model` with a `viewer viewerState` member; rename the `specViewer` field to `spec`; leave render (`vp`/`vpReady`/`renderCache`/…) and layout (`width`/`height`) flat; keep `mode`/`prevMode`/`errMsg`/`helpOpen` and the worktree overlay fields top-level
- [x] 2.3 Update `New` / `NewSinglePath` for the grouped field names (no new map init needed — `renderCache` stays flat)
- [x] 2.4 Mechanically update field accesses across `internal/ui/` (`model.go`, `update.go`, `viewer.go`, `index.go`, `worktrees.go`, `mouse.go`, `view.go`, `viewport.go`, `spec.go`, `tasks.go`): `m.tab`→`m.viewer.tab`, `m.changeIdx`→`m.viewer.changeIdx`, `m.specIdx`→`m.viewer.specIdx`, `m.specViewer`→`m.spec`. Render/layout accesses unchanged.
- [x] 2.5 Update test files' field accesses and `Model{…}` literals where tests reach into the renamed fields (`tab:`/`changeIdx:` → nested `viewer: viewerState{…}`; `m.specViewer`/`m.specIdx`/`m.tab`/`opened.tab` → renamed), keeping assertions identical (no behavior change)
- [x] 2.6 `go build ./...` and `go test ./...` green after the rename, before introducing `setMode`

## 3. setMode()

- [x] 3.1 Add `func (m *Model) setMode(next Mode)` to `model.go` (or `update.go`) per design: clear `spec.JumpTarget`/`spec.FocusMode` when leaving `ModeViewingSpec`; clamp `viewer.tab` to an available tab when entering `ModeNormal`/`ModeViewingArchive`; clamp `viewer.specIdx` against the destination change's `SpecFiles`; clamp `spec.Cursor` into `projectSpecs` when entering `ModeViewingSpec`
- [x] 3.2 Document in the `setMode` doc comment that it does NOT touch `renderCache` (cache policy stays with callers)
- [x] 3.3 Migrate every `m.mode = X` call site to `m.setMode(X)` (viewer.go, index.go, worktrees.go, config.go, model.go), preserving each site's pre-transition field assignments and `prevMode` bookkeeping
- [x] 3.4 Remove now-redundant inline cursor/focus resets in `activateIndexItem` that `setMode` subsumes, keeping only the destination-intent assignments (e.g. setting `spec.JumpTarget`/`spec.FocusMode` for `indexKindRequirement` before the transition)

## 4. Receiver convention

- [x] 4.1 Add the value/pointer receiver convention comment block atop `update.go` per design

## 5. Verification

- [x] 5.1 Grep gate: confirm the only `m.mode =` assignment outside the `setMode` body is none — all transitions route through `setMode` (only `m.mode = next` inside `setMode` remains)
- [x] 5.2 `gofmt`/`go vet ./...` clean; run the project's Go style checks (golangci-lint not installed locally; CI lint is `go vet`, which passes)
- [x] 5.3 `go test -race -cover ./...` green with the test suite behavior unmodified (only field-access renames from 2.5)
- [x] 5.4 `openspec validate model-struct-cohesion` passes
- [x] 5.5 Add `setmode_test.go` regression tests pinning the setMode invariants (leave-spec focus clear, tab/cursor clamps, renderCache preserved) — covers the new tui-viewer spec scenarios and the focus-reset consolidation in `activateIndexItem`
