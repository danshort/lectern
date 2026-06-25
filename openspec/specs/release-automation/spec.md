# release-automation Specification

## Purpose
TBD - created by archiving change automate-releases. Update Purpose after archive.
## Requirements
### Requirement: Commit-driven version and changelog
The release process SHALL derive the next version and changelog entries from Conventional Commit messages merged to `main`, rather than from manual edits. `feat:` commits SHALL produce a minor bump, `fix:` a patch bump, and a `!`/`BREAKING CHANGE` a major bump.

#### Scenario: Feature merge proposes a minor release
- **WHEN** one or more `feat:` commits are merged to `main` since the last release
- **THEN** release automation proposes a release whose version is a minor bump and whose changelog lists those features

#### Scenario: Only chore/docs merges propose no release
- **WHEN** only `chore:`/`docs:` commits are merged since the last release
- **THEN** no version bump is proposed (no release is cut for trivial changes)

### Requirement: Human-gated release via a release PR
Releases SHALL NOT be published automatically on every merge. The automation SHALL maintain a release pull request that accumulates the pending version bump and changelog; a release SHALL be created only when a maintainer merges that release PR.

#### Scenario: Release PR accumulates changes
- **WHEN** multiple feature merges land before any release
- **THEN** a single release PR is kept up to date with the combined version bump and changelog

#### Scenario: Merging the release PR cuts the release
- **WHEN** the maintainer merges the release PR
- **THEN** a version tag and GitHub release are created

### Requirement: Build and publish binaries on release
When a release is created, the pipeline SHALL build binaries for the supported platforms, attach them (with checksums) to the GitHub release without discarding the release notes, and update the Homebrew tap formula — all within the same workflow run that created the release.

#### Scenario: Binaries and tap updated on release
- **WHEN** a release tag is created by the automation
- **THEN** the GitHub release gains the platform tarballs and `checksums.txt`, and the Homebrew tap formula is updated to the new version

#### Scenario: Release notes preserved
- **WHEN** GoReleaser uploads binaries to the release created by the automation
- **THEN** the changelog notes already on the release are retained (binaries are appended, not replaced)

### Requirement: Changelog input is controlled via squash merge and conventional PR titles
Because the changelog is generated from commit subjects on the default branch, the repository SHALL squash-merge pull requests so each merged PR contributes exactly one commit whose subject is the PR title. PR titles SHALL be validated as Conventional Commits by a required check so that only well-formed, user-readable titles reach the changelog.

#### Scenario: A merged PR yields one changelog-ready commit
- **WHEN** a pull request is merged
- **THEN** it produces a single squash commit whose subject is the PR title

#### Scenario: A non-conventional PR title is rejected
- **WHEN** a pull request's title is not a valid Conventional Commit (e.g. missing a `type:` prefix)
- **THEN** the PR-title check fails, signalling the title must be fixed before merge

