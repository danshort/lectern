## 1. Implementation

- [x] 1.1 In `internal/ui/model.go`, modify the `New()` function so that when `len(project.Changes) == 0`, the model loads archive changes and project specs, builds index items, and sets `m.mode = ModeIndex` instead of defaulting to `ModeNormal`.
- [x] 1.2 Verify that when a new change appears while in the index (via the existing tick handler), the user can press Enter to open it normally without any mode-related issues.

## 2. Verification

- [x] 2.1 Run `go build ./...` to ensure the change compiles
- [x] 2.2 Run `go test ./...` to ensure existing tests pass
- [x] 2.3 Manually verify: run dossier in a project with no active changes, confirm the index view appears immediately instead of the welcome screen
