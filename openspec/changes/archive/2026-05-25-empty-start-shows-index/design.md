## Context

When dossier starts with no active changes (`len(project.Changes) == 0`), the default mode zero value (`ModeNormal`) causes `View()` to render `emptyViewContent()` — a static message. The user must press `a` or `Esc` to reach `ModeIndex`. The index view already handles the empty state (showing "No active changes" alongside specs and archives), making it a strictly better default landing experience.

The change is minimal: alter the `New()` constructor to enter `ModeIndex` immediately when no changes exist.

## Goals / Non-Goals

**Goals:**
- Start in `ModeIndex` when no active changes exist
- Preserve all existing index functionality (navigation, real-time updates, help bar)
- Keep the empty welcome screen as a fallback

**Non-Goals:**
- Removing `emptyViewContent()`
- Modifying `ModeIndex` rendering or navigation
- Changing behavior when changes DO exist at startup

## Decisions

### Decision 1: Set mode in `New()` rather than `View()` or `Update()`

**Chosen:** In `New()`, when `len(project.Changes) == 0`, initialize index data and set `m.mode = ModeIndex`.

The `View()` function already routes `ModeIndex` to `viewIndexContent()`. The `loadViewport()` method already handles `ModeIndex` by calling `refreshIndexViewport()`. Both paths work correctly as long as the index data (archive changes, project specs, index items) is initialized. No changes needed in these methods.

**Alternative:** Check in `View()` and route to index content. Rejected because it would be a view-only change — the model state would remain `ModeNormal`, causing inconsistencies in key handling (`h`/`l` would incorrectly trigger change navigation).

### Decision 2: Load archive changes and specs eagerly

**Chosen:** Load `archiveChanges` and `projectSpecs` in `New()` alongside setting the mode, matching what `enterIndex()` does.

This ensures the index renders all three sections on the first frame. Without it, the first render would show empty sections until a tick reloaded them.

### Decision 3: Keep `emptyViewContent()` unchanged

The empty welcome screen remains as a fallback. It's still reachable if the tick handler transitions to `ModeNormal` with `len(project.Changes) == 0` after the project was reloaded. Retaining it costs nothing and avoids a potential nil/blank view.

## Risks / Trade-offs

- **Slightly slower cold start with no changes:** `New()` now makes additional disk reads (archive listing, spec loading). Mitigation: these reads are fast; archives and specs are typically small directories. Users with projects that have many specs/archives but no changes were already going to pay this cost on their first `a`/`Esc` keystroke anyway.

- **Tick handler interaction:** When a new change appears while in `ModeIndex`, the existing tick handler reloads the project but keeps the mode as `ModeIndex`. The user can press Enter to open the change, or `a`/`Esc` (which re-enters index — a no-op). This is the desired behavior.
