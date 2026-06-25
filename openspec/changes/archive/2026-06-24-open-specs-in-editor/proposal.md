## Why

The `e` ("open in editor") shortcut currently works only while viewing a change's artifacts, not while viewing a project spec. A user reading a spec — or focused on a single requirement — has no way to jump into their editor to fix it and must leave the TUI to do so. Extending `e` to spec views closes that gap and makes the editing workflow consistent across the tool. (GitHub issue #18.)

## What Changes

- Handle the `e` key in `ModeViewingSpec` so the user can open the spec being viewed in `$EDITOR` (falling back to `vi`), reusing the existing editor-launch mechanism (`tea.ExecProcess`, suspend/resume, mouse re-arming).
- Make `e` available both for the full-spec view and for the single-requirement focus view; in both cases it opens the spec's `openspec/specs/<name>/spec.md` file (requirements are sections within that single file, not separate files).
- Reload the spec content on return from the editor, consistent with the existing change-artifact reload behavior.
- Add `e: edit` to the spec-view help bar (both full and focus variants).

## Capabilities

### New Capabilities
<!-- None: this extends an existing capability. -->

### Modified Capabilities
- `editor-launch`: extend the `e` shortcut beyond change artifacts to also open the current project spec from `ModeViewingSpec`, including when a single requirement is focused.

## Impact

- Code: `internal/ui/spec.go` (key handling in `ModeViewingSpec`), `internal/ui/model.go` (resolve the on-disk path of the spec being viewed), `internal/ui/view.go` (help-bar text for spec modes), and the `editorReturnMsg` reload path in `internal/ui/update.go`.
- No new dependencies, no config or CLI changes, no breaking changes.
