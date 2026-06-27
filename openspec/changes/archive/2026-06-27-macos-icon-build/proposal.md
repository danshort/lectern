## Why

The v0.20.0 release shipped the macOS app **without its icon** (placeholder in Finder/Dock). Two bugs: (1) the Icon Composer `.icon` only compiles with `actool` from **Xcode 26.4+**, but the `macos-15` CI runner has **26.3**, so the compile failed; (2) `package.sh` "degraded gracefully" and shipped an icon-less build, hiding the failure (its warning even mis-blamed "needs Xcode 26+", which the runner had). Separately, the unnotarized-launch instructions (cask caveat + README) tell users to **right-click → Open**, which **macOS 15 Sequoia removed** — so users get stuck at the Gatekeeper block.

## What Changes

- **Ship a pre-compiled icon, toolchain-independently.** Commit the compiled `Assets.car` + `lectern.icns` (regenerated from `lectern.icon` with Xcode 26.4+) under `Resources/AppIcon/`. `package.sh` embeds those, so the icon no longer depends on the CI runner's Xcode. A live `actool` compile remains a fallback.
- **Fail loudly, never silently icon-less.** If `package.sh` can embed neither the committed assets nor a freshly compiled icon, it errors and exits non-zero instead of shipping a placeholder build.
- **Add `scripts/regen-icon.sh`** to regenerate the committed assets when `lectern.icon` changes (errors clearly if the local Xcode is too old).
- **Fix the Gatekeeper instructions** in the cask `caveats` and README for macOS 15+: use `xattr -dr com.apple.quarantine …` or **System Settings → Privacy & Security → Open Anyway**; note that right-click → Open no longer bypasses on Sequoia+.

## Non-goals

- Notarization (#67) — that removes the Gatekeeper step entirely; this just corrects the interim guidance.
- Changing the icon artwork or the `.icon` source.

## Capabilities

### Modified Capabilities

- `macos-app`: the distributed build always carries its icon (toolchain-independent, build fails if it can't), and the unnotarized-launch guidance is correct for current macOS.

## Impact

- `macos/LecternApp/scripts/package.sh` — embed committed icon, fallback compile, fail-loud.
- `macos/LecternApp/scripts/regen-icon.sh` — new.
- `macos/LecternApp/Resources/AppIcon/{Assets.car,lectern.icns}` — new committed assets.
- `macos/Casks/lectern-app.rb`, `README.md` — corrected Gatekeeper instructions.
- No app/domain/corpus changes. Ships as a `fix:` → patch release that rebuilds the app (now with the icon) and republishes the cask.
