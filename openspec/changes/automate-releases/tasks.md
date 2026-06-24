## 1. release-please configuration

- [x] 1.1 Add `release-please-config.json` (`release-type: go`, `package-name: speclio`, changelog `CHANGELOG.md`)
- [x] 1.2 Add `.release-please-manifest.json` seeded at the current published version (`0.15.0`)

## 2. Workflow

- [x] 2.1 Rewrite `.github/workflows/release.yml`: trigger on push to `main`; `release-please` job + `goreleaser` job gated on `release_created`, checking out the created tag
- [x] 2.2 Grant `contents: write` and `pull-requests: write` permissions
- [x] 2.3 Set GoReleaser `release.mode: append` so binaries attach to release-please's release without replacing notes

## 3. Docs

- [x] 3.1 Update `RELEASING.md` to document the release-please flow and the one-time "Allow GitHub Actions to create and approve pull requests" repo setting

## 4. Verification

- [x] 4.1 `goreleaser check` (valid), `release-please-config.json`/manifest parse as JSON, `actionlint` clean on the workflow
- [ ] 4.2 Post-merge: confirm release-please opens a `0.16.0` release PR; merging it produces a GitHub release with binaries and an updated tap formula
