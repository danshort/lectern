## Why

Lectern's keyboard shortcuts are discoverable only through the one-line help bar, which shows a different (and truncated) subset per mode and omits several keys entirely (e.g. `1`-`4`, `s`, `h`/`l`). New users have no single place to learn what every key does on each screen. Issue #29 asks for a `?` shortcut that opens keyboard-shortcut help, grouped by screen.

## What Changes

- Add a `?` keybinding, available from every mode (`ModeNormal`, `ModeIndex`, `ModeViewingArchive`, `ModeViewingSpec`, `ModeViewingConfig`), that opens a keyboard-shortcut **help overlay**.
- Render the overlay as a centered, bordered box drawn on top of the current screen (modal-like), listing all shortcuts **grouped by screen** (Index, Change viewer, Archive viewer, Spec viewer, Config viewer, plus a Global group), each as a key/description pair.
- Dismiss the overlay with `?`, `Esc`, or `q`, returning to the exact screen and state it was opened from (no quit, no mode change).
- Suppress the overlay's `?` trigger while the index filter input is active, so `?` can still be typed into a filter query.
- Add a `?: help` affordance to the help bar across modes so the shortcut is discoverable.

## Capabilities

### New Capabilities
- `keyboard-help`: A `?`-triggered, modal-like overlay that lists all keyboard shortcuts grouped by screen, opened from any mode and dismissed back to the originating screen without side effects.

### Modified Capabilities
<!-- None: `?` is an additive, global affordance; existing modes' requirements do not change. -->

## Non-goals

- No interactive content in the overlay: it is read-only and not scrollable in this change (the shortcut list is short enough to fit; if it ever overflows, scrolling is a future addition).
- No remapping or configuration of shortcuts; the overlay documents the fixed bindings only.
- No change to any existing shortcut's behavior, and no removal of the per-mode help bar.
- No mouse interaction for opening or dismissing the overlay.
- No persistence or "don't show again" state.

## Impact

- Affected code: `internal/ui/model.go` (a flag/field to track that the overlay is open and the mode it was opened from); `internal/ui/update.go` (intercept `?` in `dispatchKey` before per-mode dispatch, and handle dismiss keys while open); `internal/ui/view.go` (compose the centered overlay over the current view using lipgloss; add the `?: help` help-bar affordance); a new `internal/ui/help.go` (the grouped shortcut catalog and overlay rendering).
- Affected specs: new `openspec/specs/keyboard-help/spec.md`.
- Tests: new UI tests for `?` opening the overlay from each mode, dismiss keys restoring the prior mode/state, the filter-input exception, and that the rendered overlay contains each screen group and representative shortcuts.
- Dependencies: none new; reuses the existing lipgloss styling and `width`/`height` already tracked from `WindowSizeMsg`.
- No data-format or external API changes.
