## Context

`AppModel` keeps a 1.5s `worktreePollTimer` (started/stopped by `updateWorktreePolling` on the `mode` didSet) that calls `pollWorktreeChanges` to re-read the surveyed worktree changes and rebuild the sidebar on a diff. Separately, the open project's `openspec/` tree is already watched by a `DirectoryWatcher` (`startWatching` → `refreshData`), and `refreshData` calls `loadWorktrees` (git enumerate + survey). So the *current* worktree's changes already refresh via FSEvents; only the *foreign* worktrees rely on the timer.

## Goals / Non-Goals

**Goals:** event-driven foreign-worktree refresh (no 1.5s wakeups), mode-scoped, capturing the worktree set on entry/reload, no spurious rebuilds, selection preserved.

**Non-Goals:** re-running `git worktree list` on file events; watching the current worktree separately (already covered); touching the TUI.

## Decisions

### D1 — Per-foreign-worktree FSEvents watchers, mode-scoped
Add `worktreeWatchers: [DirectoryWatcher]`. `updateWorktreeWatching()` (replacing `updateWorktreePolling`) rebuilds this list: when `mode == .worktrees`, create one `DirectoryWatcher` per non-bare, non-current worktree on its `openspec/` subdir (falling back to the worktree root, mirroring `startWatching`); otherwise clear the list. Called from the `mode` didSet and from `refreshData` after `loadWorktrees` (so the watcher set tracks the freshly-enumerated worktrees while in the mode). `teardown()` clears it.
- *Why exclude the current worktree:* the main project watcher already fires `refreshData` → `loadWorktrees` for it, so a separate watcher would double-refresh.

### D2 — Watcher fires → survey + diff-gated rebuild
Each watcher's callback runs `surveyWorktreeChanges()` on the main actor. That helper is the disk-survey extracted from `loadWorktrees` (loop non-bare worktrees, `loader.loadFrom`, build `worktreeChanges`/`worktreesWithProject`) — **no git**. It compares the new survey to the current `worktreeChanges`; if unchanged it returns false (no rebuild); if changed it updates state, rebuilds the sidebar, and restores the prior selection if it still exists. This preserves the poll's "no spurious rebuilds / no selection churn" behavior. `loadWorktrees` now calls `surveyWorktreeChanges()` after enumeration.
- *Re-survey vs reload-captured-set:* re-surveying via `loadFrom` also picks up changes added/removed *within* a worktree (the old poll only reloaded the captured set), a small correctness gain at no extra git cost.

## Risks / Trade-offs

- **More open event streams** (one per foreign worktree) → bounded by the number of worktrees, created only in Worktrees mode and released on exit; FSEvents is cheap at idle (the whole point).
- **Watcher-set staleness if the worktree set changes mid-mode** → re-synced whenever `refreshData`/`loadWorktrees` runs (e.g., the main watcher firing), consistent with "set captured on entry/reload."
- **Not headlessly testable** → `swift build` + manual QA; `surveyWorktreeChanges` keeps the (diffing) logic small and self-contained.

## Migration Plan

macOS-only, swaps the refresh mechanism behind the same observable behavior (mode-scoped live worktree progress). Rollback is reverting. No data/contract change.

## Open Questions

- None.
