## Why

`Model` (internal/ui/model.go) holds change-viewer, index, spec, render, and layout state as a mix of flat fields and half-grouped sub-structs. The `mode` + `tab` + `changeIdx` + `specViewer.{Cursor,FocusMode,JumpTarget}` combination can form invalid states that are prevented only by convention and re-checked ad hoc across ~14 transition sites. The mutate-and-return value-receiver pattern compounds the risk: a future handler that forgets `return m` silently drops state, with no compiler help. This is tracked tech debt (issue #37) surfaced in code review.

## What Changes

- Group the `Model` fields implicated in invalid mode/tab/cursor combinations into a cohesive sub-struct, and rename the existing spec sub-struct for symmetry:
  - `viewer` — `changeIdx`, `tab`, `specIdx` (active + archive viewing position)
  - rename the `specViewer` field to `spec`
  - `render` state (`vp`, `vpReady`, `renderCache`, …) and `layout` state (`width`, `height`) stay flat: they take no part in any invalid state and grouping them is pure cosmetic churn (~90 test-literal regroupings) with no bearing on the invariant `setMode` enforces. Scope deliberately narrowed from issue #37's literal "viewer/render/layout" wording to what serves the stated problem.
  - `mode`, `prevMode`, `errMsg`, `helpOpen` stay top-level (they are the discriminant / global chrome)
- Add `setMode(next Mode)` as the single entry point for mode transitions. It clamps the destination mode's cursors into range and resets outgoing-mode-only fields (e.g. clears `JumpTarget` and `FocusMode` when leaving `ModeViewingSpec`). It does **NOT** touch `renderCache` — cache invalidation policy stays with callers, which deliberately vary it.
- Route all `m.mode = …` assignments through `setMode()`; pull the per-field resets currently inlined in `activateIndexItem` behind it.
- Document the value/pointer receiver convention in a comment atop `update.go`.

This is a behavior-preserving refactor. No user-facing behavior changes; the existing test suite is the contract and must stay green without modification.

## Capabilities

### New Capabilities
<!-- None: this change introduces no new user-facing capability. -->

### Modified Capabilities
<!-- None: behavior-preserving internal refactor. No spec-level requirement changes. -->

## Non-goals

- No changes to user-facing behavior, key bindings, navigation, or rendering output, with **one** accepted edge-case exception: returning from the config overlay now clamps the selected tab to an available one if that tab's artifact file was deleted on disk *while config was open* (previously it left you on the now-disabled tab). Reachable only by a concurrent disk change during the overlay; strictly an improvement. See design.md Risks.
- No change to cache invalidation policy (which transitions keep vs. drop `renderCache`).
- No conversion of value-receiver handlers to pointer receivers — the Bubble Tea mutate-and-return convention is kept and documented, not replaced.
- Not splitting into stacked PRs; delivered as one change.

## Impact

- **Code**: `internal/ui/` — `model.go`, `update.go`, `viewer.go`, `index.go`, `worktrees.go`, `mouse.go`, `view.go`, `viewport.go`, `spec.go`, `tasks.go` (field-access renames for `viewer`/`spec` touch most files; `render`/`layout` accesses are untouched).
- **Tests**: `internal/ui/*_test.go` must pass unchanged. A test requiring edits signals an unintended behavior change.
- **Specs**: none — no `openspec/specs/` files change.
- **APIs / dependencies**: none.
