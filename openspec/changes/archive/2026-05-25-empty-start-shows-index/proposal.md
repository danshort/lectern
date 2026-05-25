## Why

When dossier starts with no active changes, it shows a static empty screen with a help message. This requires the user to press `a` or `Esc` to reach the index view — an unnecessary extra step. The index view already handles the empty state gracefully (showing "No active changes", specs, and archived changes), so it's a better default landing screen.

## What Changes

- On startup, when there are no active changes, dossier opens directly in `ModeIndex` instead of `ModeNormal` with an empty welcome screen
- The index shows the "Active Changes" section (empty), "Specifications", and "Archived Changes" immediately, giving the user a complete overview of the project from the first frame

## Capabilities

### New Capabilities
None.

### Modified Capabilities
- **tui-viewer**: The "Pantalla de bienvenida sin changes activos" requirement SHALL change: instead of showing a static welcome message when there are no active changes, the TUI SHALL open directly in the index view (`ModeIndex`), which already displays all project sections including the empty active changes list.

## Impact

- Affected code: `internal/ui/model.go` — `New()` function: when `len(project.Changes) == 0`, set mode to `ModeIndex` and initialize index data (archive changes, project specs, index items)
- The `emptyViewContent()` function and related view logic remain in place (still reachable if the user somehow gets into ModeNormal with no changes via edge cases)
- No breaking changes: the index view already handles the empty state, and existing keybindings (`a`/`Esc` to index, `q` to quit) work identically in `ModeIndex`

## Non-goals

- Removing the `emptyViewContent()` / welcome screen entirely (kept as fallback)
- Modifying the index view layout or behavior
- Changing how `a`/`Esc` work (they remain functional)
