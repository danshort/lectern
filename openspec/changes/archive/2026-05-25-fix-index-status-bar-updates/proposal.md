## Why

In the index view, progress bars next to active changes do not update in real-time when tasks are toggled externally. The only way to see updated progress counts is to enter a change and return to the index. This makes the index feel stale and breaks the "live" expectation established by the existing real-time refresh behaviors (which correctly detect new/deleted changes).

## What Changes

- Add in-memory task content reloading for all active changes on every tick when in `ModeIndex`, after the existing name-check passes
- Progress bars in the index now refresh within 500ms when task content changes on disk, without requiring the user to leave and re-enter the index

## Capabilities

### New Capabilities
None.

### Modified Capabilities
- **change-index**: The "Actualización en tiempo real del índice" requirement currently specifies detection of list changes (new/deleted changes). It SHALL also detect task content changes within existing changes so that progress bars update in real-time.

## Impact

- Affected code: `internal/ui/index.go` — the `handleTick()` function, specifically the `ModeIndex` branch (lines 41-45)
- No API, dependency, or configuration changes
- No breaking changes

## Non-goals

- Empty progress bars for changes without `tasks.md` (current behavior of hiding the bar when total=0 is preserved)
- Changing the tick interval (remains at 500ms)
- Reloading full project/specs/archives on every tick (only reload change task content when the name list is unchanged)
