## 1. Implementation

- [x] 1.1 In `internal/openspec/loader.go`, add a sort in `LoadFrom()` after the change-reading loop that sorts `project.Changes` by `Created` date descending (newest first). Changes without a `Created` date go last, sorted alphabetically by `Name`. Use `sort.SliceStable`.

## 2. Verification

- [x] 2.1 Run `go build ./...` to ensure the change compiles
- [x] 2.2 Run `go test ./...` to ensure existing tests pass
- [ ] 2.3 Manually verify: open dossier in a project with multiple active changes created at different times, confirm the index shows them newest-first
