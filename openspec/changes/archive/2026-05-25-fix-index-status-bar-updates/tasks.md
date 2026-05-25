## 1. Implementation

- [x] 1.1 In `internal/ui/index.go`, modify the `ModeIndex` branch of `handleTick()` so that when the name-based comparison passes (no structural changes), each active change's task content is reloaded from disk via `openspec.ReloadChange`. If any change's `Tasks.Content` or `Tasks.Present` differs from the in-memory copy, update it and mark the index as needing refresh.
- [x] 1.2 After the task-content reload loop, if any change was updated, call `buildIndexItems()` to regenerate the index item list and `refreshIndexViewport()` to re-render the viewport, preserving cursor position with the existing bounds-check logic.

## 2. Verification

- [x] 2.1 Run `go build ./...` to ensure the change compiles without errors
- [x] 2.2 Run `go test ./...` to ensure existing tests pass
- [x] 2.3 Manually verify: open dossier in index view, externally modify a `tasks.md` file in one of the active changes, confirm the progress bar updates within a few seconds without leaving the index
