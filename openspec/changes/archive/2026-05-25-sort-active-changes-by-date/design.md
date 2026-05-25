## Context

`LoadFrom()` reads change directories via `os.ReadDir`, which returns entries in filesystem order. No sorting is applied. `ListArchiveChangesFrom()` already sorts archived changes by name in reverse (newest first, since directory names are date-prefixed). The `Change.Created` field is populated from `.openspec.yaml`'s `created` field. The sort should be applied in `LoadFrom()` so all consumers (index view, tab navigation, model initialization) receive consistently ordered changes.

## Goals / Non-Goals

**Goals:**
- Sort active changes in `LoadFrom()` by `created` date, newest first
- Changes without a `created` date go last, sorted alphabetically by name as tiebreaker
- The sort applies universally (index, tab nav, header position)

**Non-Goals:**
- Sort toggle or user-configurable ordering
- Modifying the `Created` field or `.openspec.yaml` parsing

## Decisions

### Decision 1: Sort in `LoadFrom()` rather than in the UI layer

**Chosen:** Add `sort.SliceStable` after reading all changes in `LoadFrom()`.

Sorting at the data layer means all consumers benefit from consistent ordering without each having to re-sort. The index view, tab navigation (`h`/`l`), and change header position `[N/M]` all use the same slice order.

**Alternative:** Sort only in `buildIndexItems()`. Rejected because it would make the index order inconsistent with `h`/`l` navigation, and require duplicated sort logic.

### Decision 2: Sort criteria — `created` date descending, then name ascending

Changes with a valid `created` date are sorted newest first (descending). Changes without a date are placed after dated changes and sorted alphabetically by name. This matches the archive behavior (newest first) while being deterministic for undated changes.

**Tiebreaker:** If two changes have the same `created` date, they are sorted alphabetically by name (stable sort preserves input order from `os.ReadDir`).

### Decision 3: No changes to `New()` or `changeIdx`

`changeIdx` defaults to 0, which after sorting points to the newest change. This is the desired behavior. The existing order-dependent logic (checking `changeIdx` bounds) works unchanged since the slice length is the same.

## Risks / Trade-offs

- **Change order shifts on upgrade:** Users accustomed to alphabetical order will see a different order. Mitigation: this is a one-time shift to chronological order, which is more intuitive and matches archive behavior.

- **Missing `created` dates:** Older changes may lack `.openspec.yaml` or the `created` field. Mitigation: these are grouped last with alphabetical ordering, preserving determinism.
