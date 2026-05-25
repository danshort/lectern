## Context

The index view (`ModeIndex`) renders a list of active changes with progress bars showing task completion (`done/total`). The tick handler (500ms interval) currently detects only changes in the **list of change names** (new/deleted changes). If a change's `tasks.md` is modified externally — e.g., an AI agent toggles a checkbox — the in-memory `Change.Tasks.Content` is stale and the progress bar in the index does not update.

The "Actualización en tiempo real del índice" requirement in `change-index` spec mandates detection of changes on disk within 500ms, but the current implementation only checks for additions/removals, not content mutations.

## Goals / Non-Goals

**Goals:**
- Update progress bars in the index view within 500ms when `tasks.md` content changes on disk
- Preserve existing behavior for name-based change detection (new/deleted changes, new/deleted specs)
- Preserve expandedSpecs state when only task content changes (avoid resetting expansions on every tick)

**Non-Goals:**
- Reload all artifacts (proposal, design, specs) for every change on every tick
- Reload archived changes or project specs on every tick when only task content changes
- Show empty progress bars for changes without `tasks.md` (current behavior preserved)
- Change the tick interval

## Decisions

### Decision 1: Incremental task-content reload after name-check passes

**Chosen:** After the existing name comparison confirms no structural changes, iterate each active change, reload its tasks from disk via `openspec.ReloadChange`, and if any task content changed, rebuild index items and refresh the viewport.

**Alternatives considered:**

- **Always full-project reload:** Calling `openspec.LoadFrom(m.root)` on every tick. Simpler code but reads all artifacts (proposal, design, tasks, specs) and resets `expandedSpecs` on every tick, which would collapse expanded specs repeatedly. Rejected for performance and UX reasons.

- **File modification timestamp tracking:** Store `mtime` of each `tasks.md` and only reload if it changed. More complex, adds persistent state. Unnecessary given the small data size of task files and 500ms tick interval.

### Decision 2: Reload all active changes, not just the one under cursor

The progress bars for ALL active changes are visible simultaneously in the index. If only the cursor-position change were reloaded, other progress bars would remain stale. Reloading all active changes' tasks ensures the entire index view is accurate.

### Decision 3: Use existing `openspec.ReloadChange` function

`ReloadChange` already reads all four artifacts from disk. We call it per-change and only update the `Tasks` field. Other artifact fields (proposal, design, specs) are not needed for the index progress bars but will be read anyway — this is acceptable given that:
- `ReloadChange` makes at most 4 file reads per change
- Typical projects have fewer than 20 active changes
- The 500ms tick interval provides ample budget

## Risks / Trade-offs

- **Increased disk I/O on every tick:** Each tick in ModeIndex now reads `tasks.md` (and other artifact files) for all active changes, even when nothing changed. Mitigation: file reads are fast on modern SSDs; if this becomes a bottleneck, mtime-based caching can be added later.

- **`ReloadChange` reads all artifacts, not just tasks:** Avoids creating a `ReloadTasks` function, keeping the API surface small. Trade-off accepted for simplicity.

- **Cursor position preservation:** The existing cursor-adjustment logic after `buildIndexItems()` is preserved and works correctly because item count doesn't change when only task content changes.
