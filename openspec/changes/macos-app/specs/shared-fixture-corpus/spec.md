## ADDED Requirements

### Requirement: Versioned fixture corpus
The repository SHALL contain a committed corpus of OpenSpec project fixtures under `testdata/corpus/` that exercises the domain layer's non-obvious behaviors — change sort order, archive-name parsing, spec aggregation, structural validation, unreadable artifacts, and CRLF line endings.

#### Scenario: Corpus covers the drift-risk behaviors
- **WHEN** the corpus is inspected
- **THEN** it includes at least one fixture each for change ordering, malformed archive names, delta-spec validation, an unreadable artifact, and a CRLF-authored `tasks.md`

### Requirement: Golden output reproduced by every loader implementation
Each fixture SHALL have committed golden output, and every implementation of the loader (Go and any port) SHALL reproduce that output exactly when run against the corpus.

#### Scenario: Go loader matches golden
- **WHEN** the Go golden test runs the loader over the corpus and serializes the result with sorted keys
- **THEN** the output equals the committed `golden/*.json` byte-for-byte

#### Scenario: Alternate implementation matches the same golden
- **WHEN** a non-Go implementation runs against the same corpus and serializes with sorted keys
- **THEN** its output equals the same committed golden files, so divergence between implementations fails a test

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
