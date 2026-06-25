## Context

`Model` in `internal/ui/model.go` is the single Bubble Tea state struct for the TUI. It mixes already-grouped sub-structs (`index`, `specViewer`, `tasks`, `worktrees`) with flat fields for viewing position (`changeIdx`, `tab`, `specIdx`), render state (`vp`, `vpReady`, `renderCache`, `glamourRenderer`, `lastRenderWidth`, `loading`), and layout (`width`, `height`). Mode transitions (`m.mode = …`) appear at ~14 sites across `viewer.go`, `index.go`, `worktrees.go`, `config.go`, and `model.go`; each site is responsible by convention for resetting the fields the new mode requires. `activateIndexItem` is the clearest example, hand-resetting all four `specViewer` fields differently per item kind.

This is a behavior-preserving refactor. The contract is the existing test suite (`internal/ui/*_test.go`); it must pass unchanged.

## Goals / Non-Goals

**Goals**
- Make the disjoint per-mode state visible in the type by grouping flat fields into `viewer` / `render` / `layout` sub-structs.
- Funnel every mode transition through one `setMode()` chokepoint that clamps cursors and resets outgoing-mode-only state, so an invalid (mode, cursor) combination cannot be reached by a caller that forgets a reset.
- Document the value/pointer receiver convention so a future handler that forgets `return m` is caught in review.

**Non-Goals**
- No user-facing behavior change.
- No change to `renderCache` invalidation policy.
- No replacement of the value-receiver mutate-and-return pattern (it is the Bubble Tea idiom; we document it).

## Decisions

### Decision 1: Sub-struct grouping (scoped to the invalid-combo cluster)

Group only the fields implicated in the invalid mode/tab/cursor combinations the issue describes; rename the existing spec sub-struct field for symmetry. `render` and `layout` fields stay flat.

```go
type viewerState struct {
    changeIdx int  // index into project.Changes (ModeNormal)
    tab       Tab
    specIdx   int  // active spec on TabSpecs (ModeNormal + ModeViewingArchive)
}

type Model struct {
    root    string
    loader  *openspec.Loader
    project *openspec.Project

    mode, prevMode Mode   // discriminant + return target
    errMsg         string
    helpOpen       bool

    viewer    viewerState
    index     indexState
    spec      specViewerState   // field renamed from specViewer for symmetry
    tasks     taskState
    worktrees worktreesState

    // render state — flat (takes no part in any invalid combination)
    vp              viewport.Model
    vpReady         bool
    renderCache     map[Tab]string
    glamourRenderer *glamour.TermRenderer
    lastRenderWidth int
    loading         bool

    // layout state — flat
    width, height int

    projectSpecs  []openspec.ProjectSpec
    projectConfig openspec.ProjectConfig
    theme         Theme

    worktreeViewChange    openspec.Change  // archive overlay; gated by viewingWorktreeChange
    viewingWorktreeChange bool
}
```

This is a mechanical field-access rename (`m.tab` → `m.viewer.tab`, `m.changeIdx` → `m.viewer.changeIdx`, `m.specIdx` → `m.viewer.specIdx`, `m.specViewer` → `m.spec`) across the package. Render/layout accesses (`m.vp`, `m.width`, …) are untouched. The test suite gates correctness.

**Scope decision:** issue #37 literally lists `viewer/render/layout`, but the invalid combinations it explains involve only `mode`+`tab`+`changeIdx`+`specViewer`. Grouping `render`/`layout` would add ~90 cosmetic test-literal regroupings with no bearing on the invariant `setMode` enforces, so they stay flat. The grouping is scoped to what serves the stated problem.

**`specViewer` → `spec` rename**: renamed for symmetry with the other accessors (`m.index`, `m.tasks`, `m.worktrees`). Cheap, included now to avoid a half-named struct set.

### Decision 2: `setMode()` — what it owns and what it does not

`setMode` is the single entry point for changing `m.mode`. It enforces the invariants each mode owns:

```go
// setMode transitions to next, enforcing the state invariants each mode owns.
// On leaving ModeViewingSpec it clears focus-only state; on entering a mode it
// clamps that mode's cursors into range. It is the ONLY place m.mode is assigned.
//
// It deliberately does NOT touch renderCache: cache invalidation policy
// stays with callers, which vary it on purpose (tab switches keep the cache for
// instant return; moveSpec drops only TabSpecs; activateIndexItem drops all).
func (m *Model) setMode(next Mode) {
    if m.mode == ModeViewingSpec && next != ModeViewingSpec {
        m.spec.JumpTarget = ""
        m.spec.FocusMode = false
    }
    switch next {
    case ModeNormal, ModeViewingArchive:
        if !m.tabAvailable(m.viewer.tab) {
            m.viewer.tab = m.defaultTab()
        }
        // specIdx clamped against the destination change's SpecFiles
    case ModeViewingSpec:
        if m.spec.Cursor < 0 {
            m.spec.Cursor = 0
        }
        if m.spec.Cursor >= len(m.projectSpecs) {
            m.spec.Cursor = max(0, len(m.projectSpecs)-1)
        }
    }
    m.mode = next
}
```

**Scope boundary (the one subtlety):** the issue text says `setMode` should "clear caches." We deliberately exclude that. `renderCache` is kept across tab switches by design — and tab changes are not mode changes, so they would not even route through `setMode`. Folding cache-clearing in would risk a perf regression (losing instant tab return) and a behavior change. `setMode` owns **cursor clamping + cross-mode field resets** only.

**Call-site migration:** each `m.mode = X` becomes `m.setMode(X)`. Fields that a transition sets *before* the mode flip (e.g. `activateIndexItem` setting `m.spec.Cursor`/`JumpTarget`/`FocusMode` for `indexKindRequirement`) are set first, then `setMode(ModeViewingSpec)` clamps and finalizes. The destination-specific intent (focus vs. full spec view) stays at the call site; `setMode` only enforces invariants, it does not decide intent.

`prevMode` bookkeeping for `ModeViewingConfig` (`m.prevMode = m.mode` before entering config) stays at the call site, because only the caller knows it wants a return target. `setMode` does not manage `prevMode`.

### Decision 3: Receiver-convention comment

Add a comment block atop `update.go`:

> The `Update` entry point and the `dispatchKey` / `update*` handlers are **value receivers** (`func (m Model) …`): they mutate their local copy of `m` and MUST return it (`return m, cmd`). A handler that mutates `m` and forgets to return it silently drops the change — there is no compiler warning. Helpers that take a **pointer receiver** (`func (m *Model) …`, e.g. `setMode`, `enterIndex`, `mergeReloadedChange`) mutate in place and are called for their effect. Value receivers call pointer helpers freely because the receiver copy is addressable. Do not mix the two forms in one function.

## Risks / Trade-offs

- **Large rename diff** → mitigated by the test suite as the behavior contract; review focuses on `setMode` logic and the struct definition, treating field renames as mechanical.
- **`setMode` subtly changing behavior** (e.g. clamping where the old code did not) → the new clamps must match what call sites already did. Where a site did no clamp and relied on a guaranteed-valid cursor, `setMode`'s clamp is a no-op for valid input and only adds a safety floor. Verified against existing tests.
- **Missed call site** (an `m.mode = X` left un-migrated) → grep gate in tasks ensures the only assignment of `m.mode` outside `setMode` is inside `setMode` itself.
- **Accepted edge-case divergence (config return).** `handleTick` (index.go) has no `ModeViewingConfig` branch, so while the config overlay is open the normal-mode poll still runs and `mergeReloadedChange` can flip an artifact's `Present` flag if its file changes on disk. On `esc`, the new `setMode(m.prevMode)` clamps a now-unavailable `viewer.tab` to `defaultTab()`, whereas the old `m.mode = m.prevMode` left the user on the now-disabled tab. This is the **one** observable behavior change in the refactor: reachable only when a file is added/removed precisely during the config overlay, and the new behavior is strictly better (never strands the user on a disabled tab rendering absent content). Special-casing the return path to skip the clamp was rejected — it would defeat setMode's invariant. In every common flow (no disk change during config) the clamp is a no-op and the selected tab/specIdx are preserved exactly. Surfaced by adversarial verification.

## Migration Plan

Single change, no staged rollout (internal refactor, no persisted state). Land behind the existing test suite; CI green is the gate.

## Open Questions

None. Scope and `setMode` boundary settled during exploration.
