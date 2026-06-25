## 1. Shortcut catalog

- [x] 1.1 Create `internal/ui/help.go` with a declarative catalog: a slice of groups `{Title string; Shortcuts []shortcut}` where `shortcut` is `{Keys, Desc string}`, covering Global, Index, Change viewer, Archive viewer, Spec viewer, and Config viewer
- [x] 1.2 Populate each group with the bindings actually handled by that screen (cross-check `viewer.go`, `index.go`, `spec.go`, `config.go`), including `?` (help) and `q` (quit) in Global

## 2. Overlay state and dispatch

- [x] 2.1 Add a `helpOpen bool` field to `Model` in `internal/ui/model.go`
- [x] 2.2 In `dispatchKey` (`internal/ui/update.go`), when `helpOpen` is true, close on `?`/`Esc`/`q` and swallow all other keys so the underlying screen stays inert
- [x] 2.3 In `dispatchKey`, when `helpOpen` is false, open the overlay on `?` unless the index filter input is active (`m.mode == ModeIndex && m.index.FilterActive`); otherwise fall through to the existing per-mode dispatch unchanged

## 3. Overlay rendering

- [x] 3.1 Add a `renderHelpOverlay()` method in `internal/ui/help.go` that renders the catalog into a bordered box (lipgloss `RoundedBorder`, group titles via `headerStyle`, entries via existing styles)
- [x] 3.2 In `View()` (`internal/ui/view.go`), when `helpOpen` is true, center the box over the terminal with `lipgloss.Place(m.width, m.height, Center, Center, box)`, preserving the configured theme background
- [x] 3.3 Append a `?: help` affordance to the help-bar strings in `renderHelpBar` for all regular states, excluding the active-filter branch

## 4. Tests

- [x] 4.1 Test that `?` opens the overlay from each mode (`ModeNormal`, `ModeIndex`, `ModeViewingArchive`, `ModeViewingSpec`, `ModeViewingConfig`) and leaves `m.mode` unchanged
- [x] 4.2 Test that `?`, `Esc`, and `q` each dismiss the overlay, return to the originating mode, and do not quit
- [x] 4.3 Test that a non-dismiss key (e.g. `j`) is inert while the overlay is open
- [x] 4.4 Test that `?` is appended to filter text (overlay does not open) when the index filter input is active
- [x] 4.5 Test that the rendered overlay contains each screen group heading and representative shortcuts (e.g. Index `j`/`k`, `Enter`, `/`; Change viewer `1`-`4`, `h`/`l`, `Space`; Global `?`, `q`)
- [x] 4.6 Test that the help bar includes `?: help` in a regular state and omits it while the index filter input is active

## 5. Verification

- [x] 5.1 Run `gofmt`/`go vet` and `go test ./...`; build with `make build` and manually verify `?` opens/closes the overlay from each screen
