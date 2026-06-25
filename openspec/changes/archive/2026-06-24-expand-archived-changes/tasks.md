## 1. Model and state

- [x] 1.1 Add `indexKindArchivedArtifact` to the `indexItemKind` enum in `internal/ui/model.go`
- [x] 1.2 Add `ExpandedArchives map[int]bool` to `indexState` in `internal/ui/model.go`
- [x] 1.3 Initialise/reset `ExpandedArchives` wherever `ExpandedSpecs` is initialised/reset (`New`, `enterIndex`, `pollIndexMode` structural reload)

## 2. Item building

- [x] 2.1 In `buildIndexItems()`, after each archived row, append one `indexKindArchivedArtifact` sub-item per present artifact in tab order (proposal, design, specs, tasks), encoding the target `Tab` in `reqIdx`
- [x] 2.2 Add a small helper to list the present artifact tabs for a change (reusing the `Present` flags), kept consistent with `firstAvailableTab`

## 3. Rendering and hit-testing

- [x] 3.1 In `renderIndexContent()`, render `indexKindArchivedArtifact` rows indented with the artifact-type label, using the same cursor highlight style as requirement rows
- [x] 3.2 Verify `indexItemAtContentLine()` counts the new nested rows correctly (no off-by-one in the archived section)

## 4. Key and mouse handling

- [x] 4.1 Extend the `Space` handler in `updateIndex()` to toggle `ExpandedArchives` when the cursor is on an `indexKindArchived` row, rebuild items, re-anchor the cursor on the toggled archived change, and refresh the viewport
- [x] 4.2 Extend the `Enter` handler in `updateIndex()` to open `ModeViewingArchive` with `tab = Tab(item.reqIdx)` for `indexKindArchivedArtifact`
- [x] 4.3 Extend the click handler in `internal/ui/mouse.go` to select/open `indexKindArchivedArtifact` rows, mirroring the `Enter` behaviour

## 5. Filtering

- [x] 5.1 Extend `matchesFilter()` so an `indexKindArchivedArtifact` row matches when its parent archived change name or its artifact-type label matches the query

## 6. Tests

- [x] 6.1 Test that `Space` on an archived change expands it into its present artifacts in tab order, omitting absent artifacts
- [x] 6.2 Test that `Space` again collapses it and the cursor stays anchored on the archived change
- [x] 6.3 Test that `Enter` on an artifact sub-item opens `ModeViewingArchive` with the correct tab
- [x] 6.4 Test click selection/opening of an artifact sub-item (hit-testing)
- [x] 6.5 Test that filtering keeps artifact sub-items visible when the parent archived change matches
- [x] 6.6 Run `gofmt`/`go vet` and the full `go test ./...` suite

## 7. Verification

- [x] 7.1 Run `openspec validate expand-archived-changes` and confirm the change is valid
