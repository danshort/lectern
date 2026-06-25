## Why

In the index, specifications can be expanded with `Space` to preview their requirements inline, but archived changes are opaque — the only way to see what artifacts a change contains is to open it. Issue #23 asks for the same inline-expand affordance on archived changes so users can scan the documents under a change without leaving the index.

## What Changes

- Add `Space`-to-expand behaviour to archived changes in `ModeIndex`, mirroring the existing spec → requirements expansion.
- When expanded, an archived change lists the **artifact types present** for that change as nested, navigable sub-items (proposal, design, specs, tasks) in fixed tab order, showing only those that exist.
- Pressing `Enter` (or click-selecting) an artifact sub-item opens `ModeViewingArchive` focused on that artifact's tab.
- The cursor moves through artifact sub-items with `j`/`k` like requirement sub-items; expanded artifact sub-items participate in filtering by their artifact-type label.
- Real-time index refresh preserves the cursor position and expansion semantics already used for specs.

## Capabilities

### New Capabilities
<!-- None — this extends existing index behaviour. -->

### Modified Capabilities
- `change-index`: Add a requirement for expanding/collapsing archived changes with `Space`, navigating their artifact sub-items, opening an artifact sub-item with `Enter`, and the corresponding helpbar/filter behaviour.

## Non-goals

- No change to `ModeViewingArchive` itself (the read-only artifact viewer is unchanged).
- No expansion of individual spec documents within a change (artifact-type granularity only; the `specs` sub-item opens the specs tab, not per-capability documents).
- No new validation markers on archived items (archived changes remain frozen history).
- No change to active-change rows in the index (they keep their progress bars and are not expandable).

## Impact

- Affected code: `internal/ui/index.go` (item building, rendering, `Space`/`Enter` handlers, hit-testing), `internal/ui/model.go` (index item kind for archived artifacts, expansion state), `internal/ui/mouse.go` (click handling), and the index helpbar text.
- Affected specs: `openspec/specs/change-index/spec.md`.
- Tests: `internal/ui/view_test.go` and related UI tests.
- No external API, dependency, or data-format changes.
