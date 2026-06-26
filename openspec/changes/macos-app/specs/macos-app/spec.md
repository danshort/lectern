## ADDED Requirements

### Requirement: Native macOS reader for OpenSpec artifacts
The project SHALL provide a native macOS application that reads an OpenSpec project from the same `openspec/` layout the TUI reads and presents its changes and artifacts, without modifying or depending on the TUI.

#### Scenario: Open a project and browse changes
- **WHEN** the user opens a directory containing an `openspec/` folder
- **THEN** the app lists the project's active changes and lets the user navigate into each change's proposal, design, tasks, and specs

#### Scenario: TUI unaffected
- **WHEN** the macOS app is built and run
- **THEN** the Go TUI's behavior and code are unchanged

### Requirement: Faithful domain behavior via a shared contract
The app SHALL obtain changes, tasks, validation results, and worktree data through a domain layer whose behavior matches the Go implementation as enforced by the shared fixture corpus.

#### Scenario: Loader parity
- **WHEN** the app parses a project that is also present in the shared corpus
- **THEN** the parsed result matches the corpus golden output, identical to the Go loader

### Requirement: Rendered markdown with unreadable-artifact handling
The app SHALL render artifact markdown natively and SHALL surface an artifact that exists but cannot be read as a placeholder rather than as missing.

#### Scenario: Readable artifact renders
- **WHEN** an artifact file is present and readable
- **THEN** its markdown is displayed as formatted content

#### Scenario: Unreadable artifact is flagged
- **WHEN** an artifact file exists but cannot be read
- **THEN** the app shows a placeholder indicating the read failure, not an absent artifact

### Requirement: Task toggling that preserves line endings
The app SHALL let the user toggle a task checkbox in `tasks.md`, writing the change to disk while preserving the file's existing line endings.

#### Scenario: Toggle a task
- **WHEN** the user toggles a task in the app
- **THEN** the corresponding `- [ ]`/`- [x]` marker is flipped in `tasks.md` and the file's original line endings (LF or CRLF) are preserved

### Requirement: Worktrees overview
The app SHALL present the git worktrees of the project's repository, and SHALL degrade gracefully when git is unavailable or the directory is not a working tree.

#### Scenario: Worktrees listed
- **WHEN** the project is inside a git repository with multiple worktrees
- **THEN** the app lists the worktrees and marks the current one

#### Scenario: Git unavailable
- **WHEN** git is not on PATH or the directory is not a git working tree
- **THEN** the app shows an unavailable state instead of failing

### Requirement: Live reload on disk changes
The app SHALL reflect on-disk changes to the open project's `openspec/` tree without requiring a manual refresh.

#### Scenario: External edit refreshes the view
- **WHEN** a file under the open project's `openspec/` directory is modified by another process
- **THEN** the app updates its view to reflect the change
