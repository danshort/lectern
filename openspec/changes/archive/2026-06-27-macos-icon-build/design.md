## Context

`package.sh` compiled `lectern.icon` with `xcrun actool` at build time and, on failure, warned and shipped icon-less. The `.icon` format needs actool from Xcode **26.4+**; the `macos-15` runner has 26.3, so every CI build silently lost the icon (locally, on 26.6, it worked — masking the bug). `Info.plist` already carries `CFBundleIconFile`/`CFBundleIconName = lectern`, so the bundle only needs `Assets.car` + `lectern.icns` present in `Contents/Resources`.

## Decisions

### Commit the compiled icon; don't depend on the CI toolchain

The compiled outputs (`Assets.car`, `lectern.icns`) are checked into `Resources/AppIcon/`. `package.sh` copies them into the bundle — no `actool` needed at release time, so the icon is immune to the runner's Xcode version (the actual failure mode here). This trades a ~1.7 MB binary in the repo for a deterministic, portable build, which is the right call for a release artifact that was silently breaking.

`lectern.icon` stays the editable source of truth; `scripts/regen-icon.sh` recompiles it into `Resources/AppIcon/` (requires Xcode 26.4+) and is run deliberately when the art changes. `package.sh` keeps a live-compile fallback for the rare case the committed assets are absent but a capable actool is present.

### Fail loudly

The previous "degrade gracefully to icon-less" is exactly why a broken build shipped. `package.sh` now exits non-zero if it can embed neither committed nor freshly-compiled icon assets, so a missing icon fails the release (and, being a sibling job, still can't block the CLI release).

### Gatekeeper guidance for macOS 15+

macOS 15 (Sequoia) removed the right-click → Open bypass for unnotarized apps (the block dialog only offers Move to Trash / Done). The cask caveat and README now point to `xattr -dr com.apple.quarantine` (works everywhere) or System Settings → Privacy & Security → Open Anyway. This is interim; #67 notarization removes it.

## Risks / Trade-offs

- **Committed binary can drift** from `lectern.icon` if someone edits the art without running `regen-icon.sh`. Mitigated by the regen script + the live-compile fallback on capable machines and a comment in `package.sh`.
- **Repo size** grows by ~1.8 MB. Acceptable for a deterministic icon.
- **Verification** of the actual rendered icon is still visual; this change verifies the assets + `Info.plist` keys are in the bundle (done locally) and that the build fails when they're absent.
