## Why

The macOS app drives its most dynamic view — the worktrees overview — with a 1.5s repeating `Timer` (`AppModel.worktreePollTimer`), a port of the TUI's poll. But the app already uses FSEvents (`DirectoryWatcher`) for the open project's `openspec/` tree, so the native app polls on a timer for the one view that should be event-driven (#98). Switching the worktrees refresh to FSEvents removes the steady 1.5s wakeups, makes foreign-worktree progress update near-instantly, and matches how the rest of the app already watches the filesystem.

## What Changes

- Replace the worktree poll timer with **per-worktree `DirectoryWatcher`s** over each non-bare *foreign* worktree's `openspec/` tree. (The current worktree is already covered by the main project watcher, which runs `refreshData` → re-survey.)
- A worktree watcher fires → re-survey the worktree changes from disk (no `git`) and rebuild the sidebar only if something actually changed (preserving selection) — the same "no spurious rebuild" behavior the poll had.
- Watchers are **scoped to Worktrees mode**: started on entering the mode (and re-synced when the worktree set reloads), torn down on leaving the mode and on window teardown.
- Git enumeration stays on enter/reload only (unchanged); the watcher path is disk-survey only.

## Capabilities

### Modified Capabilities
- `macos-app`: the "Live worktree progress" requirement changes from periodic polling to FSEvents-driven refresh (still mode-scoped, still capturing the worktree set on entry, still no spurious rebuilds).

## Impact

- **Code:** `macos/LecternApp/Sources/LecternApp/AppModel.swift` — remove `worktreePollTimer`/`worktreePollInterval` and `start/stopWorktreePolling`/`pollWorktreeChanges`; add `worktreeWatchers: [DirectoryWatcher]`, an `updateWorktreeWatching()` lifecycle (mode-gated), a `surveyWorktreeChanges()` helper (the disk-survey extracted from `loadWorktrees`), and wire it from the `mode` didSet and after `loadWorktrees`.
- No `DirectoryWatcher` change (reused as-is). No OpenSpecKit/Go change.
- **Tests:** macOS FSEvents behavior isn't unit-testable headlessly; covered by `swift build` + manual QA. `surveyWorktreeChanges` is small and pure-ish (disk read + compare).

## Non-goals

- Re-enumerating the worktree set (`git worktree list`) on file events — kept on enter/reload, per the issue.
- Watching the current worktree separately — it's already covered by the main project watcher.
- Any change to the TUI worktrees poll (`worktrees-view` capability) — separate codebase, out of scope.
