## Why

Active changes in the index are currently listed in filesystem directory order (typically alphabetical by name), which doesn't reflect when they were created. Archived changes are already sorted newest-first by date. Sorting active changes by creation date gives a consistent chronological view across both sections and surfaces the most recent work first.

## What Changes

- Sort active changes by creation date (newest first) in `openspec.LoadFrom()`, using the `created` field from each change's `.openspec.yaml`
- Changes without a `created` date are sorted last (alphabetically by name as a tiebreaker)
- The sort applies universally — index view and tab navigation (`h`/`l` keys) both reflect the same order

## Capabilities

### New Capabilities
None.

### Modified Capabilities
- **openspec-loader**: The `LoadFrom()` function SHALL sort active changes by the `created` date from `.openspec.yaml`, newest first. Changes without a date SHALL appear after dated ones, sorted alphabetically.
- **change-index**: The "Formato de cambios activos en el índice" requirement SHALL specify that active changes are displayed in creation date order (newest first).

## Impact

- Affected code: `internal/openspec/loader.go` — `LoadFrom()` needs a sort after reading changes
- Affected code: `internal/ui/model.go` — `New()` may need to set `changeIdx = 0` explicitly if it's already 0 (zero value)
- Navigation order (`h`/`l` keys) changes from alphabetical to chronological — this is the intended behavior
- No external dependencies, no config

## Non-goals

- Adding a sort toggle for active changes (like `s` for specs)
- Changing archive change sorting (already newest-first)
- Modifying the `.openspec.yaml` schema
