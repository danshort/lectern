## MODIFIED Requirements

### Requirement: Live worktree progress
While the Worktrees mode is active, the app SHALL reflect any change in the surveyed worktree changes' task progress — in the sidebar and in an open read-only worktree artifact — without a manual reload. The app SHALL detect these changes by watching each non-bare worktree's `openspec/` tree with the filesystem watcher (FSEvents), not by polling on a timer. Watching SHALL run only while the Worktrees mode is active and SHALL stop otherwise (and on window teardown), and SHALL NOT re-enumerate the worktree set in response to file events — the set is captured on entry/reload.

#### Scenario: Foreign worktree progress updates live
- **WHEN** the user is in the Worktrees mode and a task in another worktree's change is completed externally
- **THEN** that change's progress updates in the sidebar (and in the open read-only artifact, if shown) shortly after the file is saved, without a manual reload

#### Scenario: Watching is scoped to the mode
- **WHEN** the user leaves the Worktrees mode (or closes the window)
- **THEN** the worktree watchers are stopped

#### Scenario: No spurious updates
- **WHEN** a file event resolves to no actual change in the surveyed worktree changes
- **THEN** the sidebar is not rebuilt and the current selection is preserved
