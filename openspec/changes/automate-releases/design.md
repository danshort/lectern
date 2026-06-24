## Context

The repo already releases via GoReleaser triggered by a `v*` tag push (`release.yml`), producing binaries + checksums and updating the `danshort/homebrew-tap` formula. The gap is upstream of the tag: deciding the version and writing the changelog, both manual. Commits already follow Conventional Commits loosely (`feat:`, `fix:`, `chore:`, `docs:`), which is exactly what release-please consumes.

## Goals / Non-Goals

**Goals:**
- Automated semver bump + changelog from commit history.
- A human gate before publishing (merge the release PR), not auto-publish on every merge.
- Reuse the existing GoReleaser build + tap publication unchanged.

**Non-Goals:**
- Per-merge releases; a source version file; changing the build matrix or tap.

## Decisions

- **release-please via manifest config.** `release-please-config.json` (`release-type: go`, `package-name: speclio`) + `.release-please-manifest.json` seeded at `0.15.0` (the current published version). The next merge of `feat:` commits (the already-merged #5/#6) will propose `0.16.0`.
- **Single workflow, gated build.** A tag/release created by the default `GITHUB_TOKEN` does not emit events that trigger other workflows. So instead of a separate tag-triggered build, the rewritten `release.yml` runs release-please on push to `main` and, when `release_created == true`, runs a second job that checks out the new tag and runs GoReleaser. This avoids needing a PAT purely to chain workflows.
- **GoReleaser `mode: append`.** release-please creates the GitHub release (with changelog notes) first; `append` makes GoReleaser upload binaries + checksums onto that existing release without replacing its body. (`replace` would discard release-please's notes.)
- **`GITHUB_TOKEN` for release-please, `HOMEBREW_TAP_TOKEN` for the tap.** No new secret needed for chaining. The cross-repo tap push still uses the existing `HOMEBREW_TAP_TOKEN`.

## Risks / Trade-offs

- **[Medium] Repo setting required.** release-please opens its release PR using `GITHUB_TOKEN`, which requires "Allow GitHub Actions to create and approve pull requests" (Settings → Actions → General). Without it the PR step fails. Documented in the proposal and `RELEASING.md`; a PAT is the alternative if the org forbids that setting.
- **[Low] CHANGELOG format shift.** release-please owns `CHANGELOG.md` going forward and writes its own section format; existing hand-written `v0.x` entries remain above its output. Cosmetic.
- **[Low] Can't fully e2e-test pre-merge.** The workflow only triggers on push to `main`, so the first real exercise is after this merges (release-please opens the 0.16.0 PR). Mitigated by validating config locally (`goreleaser check`, JSON parse, actionlint).
