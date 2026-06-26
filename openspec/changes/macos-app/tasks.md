## 1. Shared fixture corpus + Go golden test (Phase 1 — stands alone)

- [ ] 1.1 Create `testdata/corpus/` with fixture projects: `basic-project`, `crlf-tasks`, `unreadable-artifact`, `malformed-archive-name`, `delta-specs` (cover sort order, archive parse, spec aggregation, validation, CRLF)
- [ ] 1.2 Add `internal/openspec/golden_test.go`: walk the corpus, run the loader, `json.Marshal` with sorted keys, compare to `golden/*.json`; support a `-update` flag to regenerate
- [ ] 1.3 Add a byte-exact golden for the toggle write path: `crlf-tasks.after-toggle.tasks.md` (proves CRLF endings are preserved)
- [ ] 1.4 Add a `validation.json` golden (path → messages) covering `ValidateSpec` and `ValidateChange`
- [ ] 1.5 Generate goldens, confirm `go test ./internal/openspec/...` passes, wire into existing CI

## 2. OpenSpecKit — Swift domain port (Phase 2)

- [ ] 2.1 Scaffold `macos/OpenSpecKit/` SwiftPM package (library target + test target), no app dependency
- [ ] 2.2 Port models as `Codable` structs with `CodingKeys` matching the Go JSON field names (snake_case)
- [ ] 2.3 Define `protocol FileSystem` + `OSFileSystem`, mirroring the Go `fileSystem` interface and not-found semantics
- [ ] 2.4 Port the loader: `loadFrom`, `loadFromPath`, archive listing, `loadSpecs` (separator, absent-on-empty), config parse
- [ ] 2.5 Port tasks: `parseTasks`, `findCursorByText`, and a **separate** CRLF-safe `toggleTask` write path
- [ ] 2.6 Port validation (`validateSpec`, `validateChange`, delta rules) and `extractRequirement`
- [ ] 2.7 Port worktree porcelain parser + `normalizePath` (symlink resolution) — parser separated from `Process` invocation for testability
- [ ] 2.8 `OpenSpecKitTests`: point `Loader` at the shared `testdata/corpus/`, encode with `.sortedKeys`, assert byte-equality vs the same `golden/*.json`; assert the toggle write golden
- [ ] 2.9 Add a macOS CI lane running `swift test`; both Go and Swift golden lanes must be green

## 3. SwiftUI reader shell (Phase 3 — read-only)

- [ ] 3.1 App target `macos/LecternApp/` depending on `OpenSpecKit`; project/folder picker (security-scoped bookmark)
- [ ] 3.2 `NavigationSplitView`: sidebar of changes → artifacts; detail pane
- [ ] 3.3 Native markdown rendering of artifacts (`AttributedString` / Swift markdown lib), incl. the `⚠ couldn't read` placeholder state
- [ ] 3.4 Specs section + project config view (parity with the TUI index)

## 4. Interaction + OS integration (Phase 4)

- [ ] 4.1 Task checkbox toggle in the UI writing through the CRLF-safe `toggleTask`; integration test on a CRLF fixture
- [ ] 4.2 Worktrees view via `Process` git invocation with a 5s watchdog; graceful "unavailable" when git is absent
- [ ] 4.3 FSEvents live reload of `openspec/`; debounce; integration test that an external edit refreshes the view
- [ ] 4.4 Open-in-editor / reveal-in-Finder affordances (parity with the TUI's `$EDITOR` open)

## 5. Packaging, signing, distribution (Phase 5)

- [ ] 5.1 Code-sign the `.app` with a Developer ID; produce a `.dmg`
- [ ] 5.2 Notarize + staple in CI (signing secrets via repo secrets)
- [ ] 5.3 Publish a Homebrew **cask** alongside the existing CLI formula; document `brew install` for both
- [ ] 5.4 Extend the release workflow to build/attach the macOS artifacts on release
- [ ] 5.5 Update `README.md` with the app, screenshots, and install instructions

## 6. Verification

- [ ] 6.1 Both golden lanes green in CI (Go + Swift) on every PR
- [ ] 6.2 Manual QA matrix: browse, render, toggle (LF + CRLF files), worktrees, live reload, missing-git, unreadable-artifact
- [ ] 6.3 Confirm the TUI is byte-for-byte unchanged (no diffs under `internal/ui`, `cmd/`, `internal/openspec` logic)
- [ ] 6.4 Verify a signed+notarized build installs cleanly past Gatekeeper on a clean machine
