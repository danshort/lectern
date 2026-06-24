# Releasing

Releases are automated with [release-please](https://github.com/googleapis/release-please)
+ [GoReleaser](https://goreleaser.com/). You do **not** hand-pick versions, edit
the changelog, or push tags. The flow:

```
merge a feature PR ──▶ release-please opens/updates a
                       "chore(main): release X.Y.Z" PR
                       (version bump + CHANGELOG.md generated from commits)
        │
        ▼  (when you're ready to ship — merge that PR)
 merge the release PR ──▶ tag vX.Y.Z + GitHub release created
                      ──▶ GoReleaser builds binaries, attaches them + checksums
                      ──▶ Homebrew tap formula updated to X.Y.Z
        │
        ▼
 teammates: brew update && brew upgrade speclio
```

Both steps run in a single workflow (`.github/workflows/release.yml`) on every
push to `main`: a `release-please` job, then a `goreleaser` job gated on
`release_created`. They're in one workflow on purpose — a tag/release created by
the default `GITHUB_TOKEN` does not trigger a separate workflow.

## Conventional Commits drive the version

release-please reads commit messages since the last release:

| Commit prefix | Effect |
|---|---|
| `feat: ...` | minor bump (0.15 → 0.16) |
| `fix: ...` | patch bump (0.15.0 → 0.15.1) |
| `feat!: ...` or `BREAKING CHANGE:` in body | major bump |
| `chore:`, `docs:`, `refactor:`, `test:` | no release on their own |

So just write good commit/PR-merge messages; the version and changelog follow.

## One-time setup

These only need to be done once for the repository.

1. **Allow Actions to open PRs.** Settings → Actions → General → Workflow
   permissions → enable **"Allow GitHub Actions to create and approve pull
   requests."** Without this, release-please cannot open its release PR with the
   built-in `GITHUB_TOKEN`. (Alternative: give the workflow a PAT.)

2. **Create the Homebrew tap repo.** A public repo named `danshort/homebrew-tap`.
   GoReleaser commits the generated `Formula/speclio.rb` there on each release.

3. **Create the tap token + secret.** A
   [fine-grained PAT](https://github.com/settings/tokens) scoped to **only**
   `danshort/homebrew-tap` with **Contents: Read and write**, stored as the
   `HOMEBREW_TAP_TOKEN` secret on `danshort/speclio`.

   > Why a separate token? `secrets.GITHUB_TOKEN` can only write to the repo the
   > workflow runs in. Pushing the formula to a *different* repo needs its own token.

## Cutting a release

1. Merge feature/fix PRs to `main` as usual.
2. release-please keeps a **release PR** open with the pending version + changelog.
   Review it; merge it when you want to ship.
3. Watch the **Actions** tab — the same run builds binaries and updates the tap.
4. Verify the install:

   ```bash
   brew update
   brew upgrade speclio   # first time: brew tap danshort/tap && brew trust danshort/tap && brew install speclio
   speclio --version
   ```

That's it — no manual tagging. (To force a release with no qualifying commits,
merge a commit like `chore: release 0.16.0` or use release-please's
`Release-As:` footer.)

## Testing the build without releasing

GoReleaser can build all targets locally without tagging or publishing:

```bash
# build binaries for every platform into ./dist (no publish)
go run github.com/goreleaser/goreleaser/v2@latest build --snapshot --clean

# validate the config
go run github.com/goreleaser/goreleaser/v2@latest check
```

`./dist` is git-ignored.

## Versioning notes

- The binary reports its version via `speclio --version`, injected at build time
  through `-ldflags "-X main.version={{ .Version }}"` — do not hard-code it.
- release-please tracks the current version in `.release-please-manifest.json`;
  `release-please-config.json` configures it (`release-type: go`).
- GoReleaser uses `release.mode: append` so it adds binaries to the release that
  release-please already created, preserving the generated changelog notes.

## Notes / future work

- **Code signing (macOS):** the macOS binaries are unsigned. They install fine
  via a Homebrew *formula* (Homebrew strips the Gatekeeper quarantine attribute),
  which is why the release uses GoReleaser's `brews:` stanza rather than
  `homebrew_casks:`. If the project later moves to signed + notarized macOS
  builds, migrating to a cask becomes worthwhile.
- **Ownership:** the module path and tap currently live under the personal
  `danshort` account. If this is adopted as an internal tool, moving the repo to
  a shared org and updating the module path is the cleanest long-term home.
