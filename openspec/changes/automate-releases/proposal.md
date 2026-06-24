## Why

Releases are fully manual today: hand-write a CHANGELOG entry, pick a version, tag, push. That's easy to forget — especially during team handoff — and the version bump is a judgment call each time. Merging features to `main` does nothing user-visible until someone remembers to cut a tag, so `brew upgrade` silently stays behind `main`.

## What Changes

- Adopt [release-please](https://github.com/googleapis/release-please) to automate versioning and the changelog from Conventional Commit messages (the repo already uses `feat:`/`fix:`/`chore:`/`docs:`).
- On each merge to `main`, release-please maintains a standing "release PR" that bumps the version and updates `CHANGELOG.md`. Merging that PR creates the tag + GitHub release.
- Rework the release workflow into a single workflow: release-please runs on push to `main`; when it cuts a release, a gated GoReleaser job builds the binaries and updates the Homebrew tap. (A tag created by the default `GITHUB_TOKEN` cannot trigger a separate workflow, so the build must run in the same workflow run.)
- GoReleaser switches to `release.mode: append` so it attaches binaries to the release release-please created without clobbering its notes.
- Update `RELEASING.md` to document the new flow and the one-time repo setting it requires.

## Non-goals

- Releasing on every merge. Releases happen only when the maintainer merges the release PR — trivial/docs merges accumulate without cutting a release.
- A version file in source; release-please tracks the version via `.release-please-manifest.json` + tags.
- Changing the Homebrew tap mechanism or the GoReleaser build matrix.

## Capabilities

### New Capabilities

- `release-automation`: automated, commit-driven versioning, changelog, tagging, and binary/tap publication.

## Impact

- `.github/workflows/release.yml` — rewritten: `release-please` job + gated `goreleaser` job, triggered on push to `main`
- `release-please-config.json`, `.release-please-manifest.json` — new (release-type `go`, seeded at 0.15.0)
- `.goreleaser.yaml` — `release.mode: replace` → `append`
- `RELEASING.md` — documents the new flow
- **One-time repo setting:** "Allow GitHub Actions to create and approve pull requests" must be enabled (Settings → Actions → General) so release-please can open its PR with `GITHUB_TOKEN`
