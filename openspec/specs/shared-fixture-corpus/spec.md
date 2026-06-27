# shared-fixture-corpus Specification

## Purpose
TBD - created by archiving change macos-app. Update Purpose after archive.
## Requirements
### Requirement: Versioned fixture corpus
The repository SHALL contain a committed corpus of OpenSpec project fixtures under `testdata/corpus/` that exercises the domain layer's non-obvious behaviors — change sort order, archive-name parsing, spec aggregation, structural validation, unreadable artifacts, and CRLF line endings.

#### Scenario: Corpus covers the drift-risk behaviors
- **WHEN** the corpus is inspected
- **THEN** it includes at least one fixture each for change ordering, malformed archive names, delta-spec validation, an unreadable artifact, and a CRLF-authored `tasks.md`

### Requirement: Golden output reproduced by every loader implementation
Each fixture SHALL have committed golden output covering every public entry point — not only the parsed `Project`, but task parsing, requirement extraction, worktree-porcelain parsing, config-to-markdown, and validation — and every implementation (Go and any port) SHALL reproduce that output exactly when run against the corpus, under a written serialization contract (no `omitempty`; absent → `null`; empty slice → `[]`; empty map → `{}`; sorted keys; snake_case field names).

#### Scenario: Go loader matches golden
- **WHEN** the Go golden test runs every entry point over the corpus and serializes with sorted keys under the contract
- **THEN** the output equals the committed golden files byte-for-byte

#### Scenario: Alternate implementation matches the same golden
- **WHEN** a non-Go implementation runs against the same corpus and serializes under the same contract
- **THEN** its output equals the same committed golden files, so divergence between implementations fails a test

#### Scenario: Entry points beyond the Project tree are pinned
- **WHEN** the goldens are inspected
- **THEN** they include task parsing (with line numbers), requirement extraction, worktree-porcelain parsing, config-to-markdown, project-specs (`openspec/specs/`), and validation output

### Requirement: Cross-language-stable error representation
Because an unreadable artifact embeds an OS- and locale-specific error string, the golden for that case SHALL normalize the error — recording presence and a read-error flag with a prefix-only content match — so the golden is reproducible across languages rather than pinned to a Go runtime error string.

#### Scenario: Unreadable-artifact golden is language-stable
- **WHEN** the unreadable-artifact fixture is evaluated by any implementation
- **THEN** the golden matches on presence, the read-error flag, and the placeholder prefix, not on the raw OS error text

### Requirement: Behavior changes require fixture updates
The corpus SHALL be treated as the canonical specification of loader behavior; any change to the loader, task, validation, or worktree-parsing behavior SHALL add or modify a fixture and its golden, and the project SHALL document that golden coverage does not imply behavioral completeness (cross-platform filesystem, Unicode, YAML, and regex-dialect differences are not pinnable by the corpus).

#### Scenario: A behavior change without a fixture is rejected
- **WHEN** a pull request changes loader/task/validation/worktree-parsing behavior without adding or updating a fixture
- **THEN** the change is flagged for a missing fixture in review

### Requirement: Byte-exact task-toggle golden
The corpus SHALL pin the task-toggle write path with a byte-exact expected file so that line endings (including CRLF) are preserved on write.

#### Scenario: Toggling a CRLF task preserves CRLF
- **WHEN** a task in a CRLF-authored `tasks.md` fixture is toggled
- **THEN** the resulting file equals the committed post-toggle golden, with CRLF line endings intact

### Requirement: Both lanes enforced in CI
The continuous integration pipeline SHALL run the golden checks for every loader implementation so a one-sided behavior change is caught as a failing build.

#### Scenario: One-sided change fails CI
- **WHEN** a parsing rule is changed in one implementation without updating the golden and the other implementation
- **THEN** at least one CI lane fails

