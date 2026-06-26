## MODIFIED Requirements

### Requirement: Worktrees overview
The app SHALL present the git worktrees of the project's repository, the current worktree first and marked, and for each non-bare worktree SHALL surface its active changes with a task-progress indicator (completed / total). Selecting a change under a worktree SHALL open it read-only. The view SHALL degrade gracefully when git is unavailable, a worktree is bare, or a worktree has no `openspec/` project.

#### Scenario: Worktrees listed
- **WHEN** the project is inside a git repository with multiple worktrees
- **THEN** the app lists the worktrees, current first and marked

#### Scenario: Per-worktree active changes and progress
- **WHEN** a non-bare worktree has active changes
- **THEN** each change is listed under that worktree with a completed/total task-progress indicator

#### Scenario: Open a worktree's change read-only
- **WHEN** the user selects a change under a worktree
- **THEN** the app renders that change's artifacts read-only (no task toggling)

#### Scenario: Worktree without an OpenSpec project
- **WHEN** a worktree is bare or has no `openspec/` directory
- **THEN** it is listed with no changes, not as an error

#### Scenario: Git unavailable
- **WHEN** git is not on PATH or the directory is not a git working tree
- **THEN** the app shows an unavailable state instead of failing
