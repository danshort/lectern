# Releasing

`dossier` is released with [GoReleaser](https://goreleaser.com/). Pushing a
`vX.Y.Z` tag to `main` triggers `.github/workflows/release.yml`, which builds
binaries for Linux and macOS (amd64 + arm64), publishes a GitHub Release with
checksums, and updates the Homebrew tap.

## One-time setup

These only need to be done once for the repository.

1. **Create the Homebrew tap repo.** Create a public repo named
   `danshort/homebrew-tap`. GoReleaser commits the generated formula
   (`Formula/dossier.rb`) there on every release. It does not need any contents
   beyond what GoReleaser writes.

2. **Create the tap token.** Generate a GitHub
   [fine-grained personal access token](https://github.com/settings/tokens)
   scoped to **only** the `danshort/homebrew-tap` repo with **Contents: Read
   and write** permission.

   > Why a separate token? The built-in `secrets.GITHUB_TOKEN` can only write to
   > the repo the workflow runs in (`danshort/dossier`). Pushing the formula to
   > a *different* repo (`homebrew-tap`) requires its own token.

3. **Add the token as a secret.** In `danshort/dossier` →
   Settings → Secrets and variables → Actions, add a secret named
   `HOMEBREW_TAP_TOKEN` with the token value.

## Cutting a release

1. Make sure `main` is green (CI passes) and `CHANGELOG.md` is updated.

2. Tag and push. Use [semantic versioning](https://semver.org/):

   ```bash
   git checkout main && git pull --ff-only
   git tag v0.11.0
   git push origin v0.11.0
   ```

3. Watch the release workflow under the repo's **Actions** tab. On success it
   produces:
   - a GitHub Release at `v0.11.0` with `dossier-<os>-<arch>.tar.gz` archives
     and `checksums.txt`,
   - an updated `Formula/dossier.rb` in `danshort/homebrew-tap`.

4. Verify the install:

   ```bash
   brew update
   brew tap danshort/tap
   brew install dossier   # or: brew upgrade dossier
   dossier --version
   ```

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

- The binary reports its version via `dossier --version`. The value is injected
  at build time through `-ldflags "-X main.version={{ .Version }}"` — do not
  hard-code it.
- `release.mode: replace` means re-pushing an existing tag replaces that
  release's artifacts. Prefer cutting a new patch version over re-tagging.

## Notes / future work

- **Code signing (macOS):** the macOS binaries are currently unsigned. They
  install fine via a Homebrew *formula* (Homebrew strips the Gatekeeper
  quarantine attribute), which is why the release uses GoReleaser's `brews:`
  stanza rather than `homebrew_casks:`. If the project later moves to signed +
  notarized macOS builds, migrating to a cask becomes worthwhile.
- **Ownership:** the module path and tap currently live under the personal
  `danshort` account. If this is adopted as an internal tool, moving the repo to
  a shared org and updating the module path is the cleanest long-term home.
