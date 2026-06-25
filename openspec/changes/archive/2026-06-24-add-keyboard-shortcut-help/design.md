## Context

Lectern is a Bubble Tea (Elm-architecture) TUI. Key presses are dispatched in `internal/ui/update.go` via `dispatchKey`, which switches on `m.mode` and delegates to a per-mode handler (`updateViewer`, `updateIndex`, `updateSpec`, `updateConfig`). The view is composed in `internal/ui/view.go` as a single full-screen box (`mainViewContent` / `viewContentWithChrome`) with a one-line help bar rendered by `renderHelpBar`. There is currently **no modal/overlay pattern**: each mode owns the whole screen.

Shortcuts today are only partially advertised in the per-mode help bar, and several keys (`1`-`4`, `s`, `h`/`l`, `Space`) are undocumented for new users. Issue #29 asks for a `?` overlay that lists all shortcuts grouped by screen.

The model already tracks `width`/`height` from `WindowSizeMsg` and a `prevMode` field (used to return from `ModeViewingConfig`), and styling uses lipgloss with a shared style set in `styles.go`.

## Goals / Non-Goals

**Goals:**
- A `?` key that opens a centered, modal-like overlay listing every shortcut, grouped by screen, available from all five modes.
- Dismiss with `?`/`Esc`/`q`, returning to the exact originating screen and state with no side effects.
- A single source of truth for the shortcut catalog so the overlay stays consistent.
- Discoverability via a `?: help` help-bar affordance.

**Non-Goals:**
- No scrolling, interactivity, mouse handling, remapping, or persistence in the overlay.
- No change to existing shortcut behavior or to the per-mode help bar contents (beyond appending the `?` affordance).

## Decisions

### Overlay state: a boolean flag, not a new Mode
Add `helpOpen bool` to `Model` rather than a `ModeHelpHelp` constant. The overlay is a transient layer over the current screen, not a screen of its own — modeling it as a flag keeps the existing per-mode state (cursor, filter, tab, viewport) untouched, so "restore the originating screen" is automatic. The originating mode is already preserved because `m.mode` is never changed while the overlay is open.

Alternative considered: a dedicated `ModeHelp` with `prevMode` tracking (mirrors `ModeViewingConfig`). Rejected because it would require saving/restoring the full per-mode state of whichever screen we left, and `prevMode` is already overloaded for the config view. A flag avoids that entirely.

### Dispatch: intercept `?` and dismiss keys at the top of `dispatchKey`
Handle the overlay centrally in `dispatchKey` before the per-mode switch:
1. If `m.helpOpen`: consume `?`, `Esc`, `q` to close (set `helpOpen = false`); swallow all other keys (return `m, nil`) so the underlying screen is inert.
2. Else if the key is `?` **and** we are not capturing filter text (`!(m.mode == ModeIndex && m.index.FilterActive)`): set `helpOpen = true`.
3. Else fall through to the existing per-mode dispatch unchanged.

This keeps every per-mode `update*` function free of overlay logic and guarantees uniform behavior across modes. The filter-input guard reuses the existing `m.index.FilterActive` flag so `?` remains typeable into a filter query.

Alternative considered: adding a `?` case to each per-mode handler. Rejected — five duplicated handlers, easy to drift, and harder to keep the "swallow other keys" semantics consistent.

### Catalog: one declarative table in a new `internal/ui/help.go`
Define the shortcut catalog as a slice of groups, each `{Title string; Shortcuts []{Keys, Desc string}}`, covering Global, Index, Change viewer, Archive viewer, Spec viewer, and Config viewer. Rendering iterates the catalog, so adding a future shortcut means editing one table. This also makes a test that asserts "every group heading and representative shortcut is present" straightforward.

The catalog is hand-maintained (not auto-derived from the dispatch code, which is plain `switch` statements over `msg.String()` with no registry). A design risk noted below covers drift; a test pins the key cases.

### Rendering: lipgloss box + `lipgloss.Place` over the base view
`View()` renders the base screen as today. When `helpOpen`, render the catalog into a bordered box (`lipgloss.RoundedBorder()`, reusing `headerStyle` for group titles and `helpStyle`/default for entries), then center it with `lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)` and return that as the view content. Because the catalog is small and fixed, this fits comfortably; if the terminal is too small, the box is clamped to the available width/height.

Alternative considered: true compositing of the box over the dimmed base view. Rejected as unnecessary complexity for v1 — replacing the content with the centered box is simpler and the base screen is fully restored on dismiss.

### Help-bar affordance
Append `?: help` to the help-bar strings for the regular states in `renderHelpBar`, except the active-filter branch (which shows the live filter prompt). This is the only edit to existing help-bar text.

## Risks / Trade-offs

- **Catalog drift from real bindings** → A unit test asserts the catalog contains the key shortcuts each screen actually handles; reviewers update the table alongside any binding change. The catalog lives next to the UI package so it is easy to find.
- **Small terminals clip the box** → `lipgloss.Place` centers within `width`/`height`; the box border/padding are minimal and content is short, so realistic terminal sizes fit. Extreme sizes degrade gracefully (clipped) rather than crashing.
- **`q` closing the overlay differs from `q` quitting elsewhere** → Intentional and documented in the spec; while the overlay is open, `q` is a dismiss key, not quit, so a user who opened help can always back out without leaving the app.
- **Background color / view chrome** → The overlay reuses the configured theme background via the existing `tea.View` background handling, so it stays visually consistent with the rest of the app.

## Open Questions

- None blocking. Whether to make the overlay scrollable can be revisited if the catalog grows beyond a typical terminal height (out of scope here).
