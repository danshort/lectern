## Context

The `e` shortcut is implemented inline in `internal/ui/viewer.go` (the `case "e":` block, ~lines 114-132). It resolves the active change artifact via `Model.artifactPath()`, splits `$EDITOR` (fallback `vi`), and launches it with `tea.ExecProcess`, posting an `editorReturnMsg` on exit. The `editorReturnMsg` handler in `internal/ui/update.go` reloads the current change and re-renders.

Project specs are a separate concern: they are loaded into `m.projectSpecs []openspec.ProjectSpec` via `loader.LoadProjectSpecsFrom(root)`, and `ModeViewingSpec` renders from the in-memory `m.projectSpecs[m.specViewer.Cursor].Content` (see `loadViewportForSpec` in `internal/ui/viewport.go`). Key handling for that mode lives in `internal/ui/spec.go` (`updateSpec`), which today has no `e` case. Requirements are `### Requirement:` sections inside a single `spec.md`; focus mode just renders one extracted block.

## Goals / Non-Goals

**Goals:**
- Pressing `e` in `ModeViewingSpec` opens the viewed spec's `openspec/specs/<name>/spec.md` in `$EDITOR`, for both full-spec and requirement-focus views.
- On return, the spec content is reloaded from disk so external edits show immediately, staying in `ModeViewingSpec` (full vs. focus preserved).
- Editor-launch logic is shared between the change viewer and the spec viewer rather than duplicated.

**Non-Goals:**
- Jumping the editor to a specific requirement's line/offset (always opens the spec file; requirements are sections of one file).
- Adding `e` to `ModeIndex` (the change-side `e` is viewer-only; we keep the spec behavior consistent with that).
- Per-requirement files or any change to on-disk spec structure.

## Decisions

**1. Extract a shared `openInEditor(path string) tea.Cmd` helper.**
Move the `$EDITOR`-splitting / `vi`-fallback / `exec.Command` / `tea.ExecProcess` logic (with its `#nosec` annotation) out of the inline `viewer.go` block into one method on `*Model` in a shared location (e.g. `viewer.go` or a small `editor.go`). Both the change viewer and `updateSpec` call it. *Why over copy-paste:* a second literal copy of the security-sensitive exec block invites drift; one helper keeps the `#nosec` rationale and parsing in a single place.

**2. Resolve the spec path with a small `currentSpecPath()` method.**
Add a `*Model` method returning `filepath.Join(m.root, "openspec", "specs", m.projectSpecs[m.specViewer.Cursor].Name, "spec.md")`, or `""` when the cursor is out of range. This mirrors the existing `artifactPath()` pattern and keeps `updateSpec` thin. The same path serves full and focus mode.

**3. Handle `e` in `updateSpec`.**
Add `case "e":` to `internal/ui/spec.go`: resolve `currentSpecPath()`; if non-empty, return `m, m.openInEditor(path)`. No focus-mode branching needed — the file is identical.

**4. Reload project specs on return when in `ModeViewingSpec`.**
Extend the `editorReturnMsg` handler in `update.go`: when `m.mode == ModeViewingSpec`, refresh `m.projectSpecs` via `m.loader.LoadProjectSpecsFrom(m.root)` (clamping `m.specViewer.Cursor` if the list shrank) before calling `m.loadViewport()`. The existing change-reload branch stays gated to change modes so spec edits don't trigger needless change reloads.

**5. Advertise `e: edit` in the spec help bar.**
Update `renderHelpBar()` in `internal/ui/view.go` for both `ModeViewingSpec` variants (full: `j/k: scroll  Esc: index  q: quit`; focus: `h/l: ...  j/k: scroll  Esc: index  q: quit`) to include `e: edit`.

## Risks / Trade-offs

- [Spec list changes during edit (e.g. file removed)] → After reload, clamp `m.specViewer.Cursor`; `loadViewportForSpec` already renders `(spec not available)` for out-of-range cursors, so the worst case is a graceful placeholder.
- [Mouse tracking after editor return] → No new risk: reusing the same `tea.ExecProcess` path means mouse mode is re-applied on the next render frame exactly as it is for change artifacts.
- [Refactor touches the security-sensitive exec block] → Mitigated by moving it verbatim (including the `#nosec G204 G702` comment) into the helper and covering both call sites with tests.
