## Context

The index (`ModeIndex`) already supports expanding a spec with `Space` to reveal its requirements inline. The machinery for this lives in `internal/ui/index.go`:

- `indexState.ExpandedSpecs map[int]bool` tracks which specs are expanded (`model.go`).
- `buildIndexItems()` flattens the three sections into a single `[]indexItem` slice, inserting `indexKindRequirement` rows after an expanded `indexKindSpec` row.
- `renderIndexContent()` renders the flat list and computes the cursor line; `indexItemAtContentLine()` does the inverse for mouse hit-testing.
- The `Space` handler in `updateIndex()` toggles `ExpandedSpecs`, rebuilds items, and re-anchors the cursor on the toggled spec.

Archived changes (`indexKindArchived`) are currently leaf rows: `Enter` opens `ModeViewingArchive`. An archived change is an `openspec.Change` whose `Proposal`, `Design`, `Specs`, and `Tasks` fields each carry a `Present` flag. We want archived changes to expand into their present artifact types, mirroring the spec → requirements flow.

## Goals / Non-Goals

**Goals:**
- `Space` toggles expansion of the archived change under the cursor, listing its present artifact types (proposal, design, specs, tasks) in tab order as nested sub-items.
- `Enter`/click on an artifact sub-item opens `ModeViewingArchive` on that artifact's tab.
- Expansion integrates with existing cursor re-anchoring, filtering, real-time refresh, and mouse hit-testing without regressing spec expansion.

**Non-Goals:**
- Changing `ModeViewingArchive` rendering or keybindings.
- Expanding `specs` into per-capability documents (the `specs` sub-item maps to the specs tab).
- Adding validation markers to archived items.

## Decisions

### Decision: New `indexKindArchivedArtifact` item kind
Add a fourth nested item kind alongside `indexKindRequirement`. The existing `indexItem` struct already has the fields needed: reuse `idx` for the archived-change index (into `ArchiveChanges`) and `reqIdx` to carry the target `Tab`. This avoids growing the struct and mirrors how `indexKindRequirement` reuses `idx`+`reqIdx`.

- *Alternative considered:* a separate `tab Tab` field on `indexItem`. Rejected — `reqIdx` is already a free per-row integer used only by nested rows, so reusing it keeps the struct minimal and the two nested kinds symmetric.

### Decision: Track expansion in a new `ExpandedArchives map[int]bool`
Keep archived-change expansion state separate from `ExpandedSpecs` (keyed by archived-change index). Reset it in the same places `ExpandedSpecs` is reset (`enterIndex`, `pollIndexMode` structural reload, `New`).

- *Alternative considered:* one shared map. Rejected — spec indices and archived indices are different namespaces; a shared map would alias unrelated rows.

### Decision: Sub-items are present artifacts in fixed tab order
`buildIndexItems()` appends artifact sub-items after an expanded archived row by iterating `TabProposal, TabDesign, TabSpecs, TabTasks` and emitting a row only when the corresponding `Present` flag is true. This reuses the existing `tabAvailable`-style presence checks and yields a stable, predictable order. The sub-item label is the artifact type name (the existing `tabLabels` value).

- *Alternative considered:* listing every spec file under `specs/`. Rejected per proposal Non-goals — artifact-type granularity keeps the mirror with spec→requirements clean and the row count bounded.

### Decision: `Enter`/click opens `ModeViewingArchive` on the chosen tab
The `Enter` handler gains a branch for `indexKindArchivedArtifact` that sets `ArchiveCursor = item.idx`, `tab = Tab(item.reqIdx)`, and `mode = ModeViewingArchive` — the same target as selecting the parent archived row, but with the tab pre-selected instead of `firstAvailableTab`. Mouse selection in `mouse.go` mirrors this.

### Decision: Render, hit-test, filter, and helpbar updates
- `renderIndexContent()`: the archived-section loop already walks `Items[specEnd:]`; extend it to render `indexKindArchivedArtifact` rows with the same indented, cursor-aware style used for requirement rows.
- `indexItemAtContentLine()`: the archived-section loop counts one content line per visible item, so nested artifact rows are counted automatically — verify no off-by-one against the new rows.
- `matchesFilter()`: an `indexKindArchivedArtifact` row matches when the parent change name or the artifact label contains the query, so expanded rows don't vanish while their parent matches.
- Helpbar already advertises `Space: expand`; no text change needed, but confirm the wording still reads correctly now that both specs and archived changes expand.

## Risks / Trade-offs

- [Cursor re-anchoring after toggle drifts] → Reuse the exact re-anchor loop from the spec `Space` handler, matched on `kind == indexKindArchived && idx == archIdx`, covered by a test.
- [Mouse hit-testing off-by-one from new nested rows] → The archived loop already iterates `Items` by visibility; add a click test that selects an artifact sub-item to lock the behaviour.
- [Real-time reload collapses user's expansion] → Acceptable and consistent with current spec behaviour: a structural disk change resets expansion maps; non-structural task refreshes preserve them.

## Open Questions

None — the granularity question (artifact types vs. individual documents) was resolved to artifact types.
