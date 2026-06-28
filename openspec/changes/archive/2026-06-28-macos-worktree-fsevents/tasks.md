## 1. Survey helper

- [x] 1.1 Extract the disk-survey from `loadWorktrees` into `surveyWorktreeChanges() -> Bool` (loop non-bare worktrees, `loadFrom`, build `worktreeChanges`/`worktreesWithProject`); return whether the survey differs from current state. `loadWorktrees` calls it after git enumeration.

## 2. FSEvents watchers (replace the timer)

- [x] 2.1 Remove `worktreePollTimer`, `worktreePollInterval`, `startWorktreePolling`, `stopWorktreePolling`, `pollWorktreeChanges`
- [x] 2.2 Add `worktreeWatchers: [DirectoryWatcher]` and `updateWorktreeWatching()`: in Worktrees mode, one watcher per non-bare, non-current worktree on its `openspec/` dir (root fallback); else clear. On fire → `surveyWorktreeChanges()` and, if changed, rebuild the sidebar preserving selection
- [x] 2.3 Call `updateWorktreeWatching()` from the `mode` didSet (replacing `updateWorktreePolling`) and from `refreshData` after `loadWorktrees`; clear watchers in `teardown()`

## 3. Verification

- [x] 3.1 `swift build` (LecternApp) + `swift build`/`swift test` (OpenSpecKit) clean; golden green
- [ ] 3.2 Manual: open Worktrees mode; complete a task in a sibling worktree's change on disk → its progress updates shortly after save with no manual reload
- [ ] 3.3 Manual: leave Worktrees mode → confirm no further worktree refreshes (watchers released); re-enter works
