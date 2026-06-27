## ADDED Requirements

### Requirement: App icon
The packaged `.app` SHALL bundle the project's app icon, compiled from the Icon Composer source (`lectern.icon`) into the bundle and referenced by `Info.plist`, so the app shows proper branding in Finder, the Dock, and About. When the build toolchain cannot compile the icon format, packaging SHALL warn and still produce a working (icon-less) build rather than fail.

#### Scenario: Packaged app includes its icon
- **WHEN** the app is packaged on a toolchain that supports the icon format
- **THEN** the bundle contains the compiled icon (`Assets.car` + `lectern.icns`) and `Info.plist` references it via `CFBundleIconName`

#### Scenario: Older toolchain degrades gracefully
- **WHEN** the build toolchain cannot compile the `.icon`
- **THEN** packaging warns and still produces a working build without the icon, rather than failing
