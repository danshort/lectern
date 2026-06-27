## Why

The macOS app ships with the generic placeholder icon. An Icon Composer icon
(`lectern.icon`) now exists; the packaged app should use it so it has proper
branding in Finder, the Dock, and About.

## What Changes

- Compile the Icon Composer `macos/LecternApp/lectern.icon` into the `.app`
  bundle during packaging (`actool` → `Assets.car` + `lectern.icns`) and
  reference it from `Info.plist` (`CFBundleIconName`/`CFBundleIconFile`).
- Packaging degrades gracefully when the toolchain can't compile the `.icon`
  (Xcode < 26): warn and ship without the icon rather than fail.

## Non-goals

- A toolbar/in-app icon or any UI change; this is the bundle app icon only.

## Capabilities

### Modified Capabilities

- `macos-app`: the packaged app bundles its app icon.

## Impact

- `macos/LecternApp/lectern.icon` (icon source), `Resources/Info.plist` (icon
  keys), `scripts/package.sh` (compile step), `.github/workflows/macos-app-release.yml`
  (select an icon-capable Xcode). No code / domain / corpus changes.
