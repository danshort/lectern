# Tasks

## 1. Toolchain-independent icon
- [x] 1.1 Commit pre-compiled `Resources/AppIcon/{Assets.car,lectern.icns}` (regenerated from `lectern.icon` on Xcode 26.4+)
- [x] 1.2 `package.sh` embeds the committed assets; live `actool` compile is only a fallback
- [x] 1.3 `package.sh` fails (non-zero) if it can embed neither committed nor freshly compiled icon
- [x] 1.4 Add `scripts/regen-icon.sh` to regenerate the committed assets (errors if Xcode too old)

## 2. Gatekeeper instructions for macOS 15+
- [x] 2.1 Cask `caveats`: use `xattr` / System Settings → Open Anyway; note right-click → Open is gone on Sequoia+
- [x] 2.2 README first-launch note: same correction

## 3. Verify
- [x] 3.1 `package.sh` builds with the icon embedded (`Assets.car` + `lectern.icns` in the bundle, `Info.plist` references it)
- [ ] 3.2 Post-release: after the patch release, `brew upgrade` shows the app icon (no placeholder)
