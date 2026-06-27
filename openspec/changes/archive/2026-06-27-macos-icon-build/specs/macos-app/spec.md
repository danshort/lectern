## MODIFIED Requirements

### Requirement: App icon
The packaged `.app` SHALL bundle the project's app icon (`Assets.car` + `lectern.icns`, referenced by `Info.plist` via `CFBundleIconName`/`CFBundleIconFile`) so the app shows proper branding in Finder, the Dock, and About. The icon SHALL be embedded independently of the build toolchain by shipping pre-compiled assets (regenerated from the `lectern.icon` source when it changes). If packaging can embed neither the committed pre-compiled assets nor a freshly compiled icon, it SHALL fail rather than produce an icon-less build.

#### Scenario: Packaged app includes its icon
- **WHEN** the app is packaged
- **THEN** the bundle contains the compiled icon (`Assets.car` + `lectern.icns`) and `Info.plist` references it

#### Scenario: Icon does not depend on the build toolchain
- **WHEN** the app is packaged on a runner whose toolchain cannot compile the `.icon` format
- **THEN** the committed pre-compiled icon assets are used, so the build still carries the icon

#### Scenario: Missing icon fails the build
- **WHEN** packaging can embed neither the committed assets nor a freshly compiled icon
- **THEN** packaging fails with an error rather than silently shipping an icon-less build

### Requirement: Distributable signed build
The app SHALL be packaged as a `.app` bundle distributed via a Homebrew **cask** alongside the CLI. Signing is phased: until a Developer-ID certificate is available it SHALL be at least **ad-hoc** code-signed (required to run on Apple Silicon); once available it SHALL be **Developer-ID** signed (hardened runtime) and **notarized**. While the distributed build is unnotarized, the install instructions SHALL document a one-time step to open it past Gatekeeper that is valid on current macOS; that guidance SHALL be removed once notarized builds ship.

#### Scenario: Build runs on Apple Silicon
- **WHEN** the app is packaged
- **THEN** the resulting `.app` carries at least an ad-hoc code signature and launches on Apple Silicon

#### Scenario: First-launch Gatekeeper guidance while unnotarized
- **WHEN** a user installs an unnotarized build
- **THEN** the install instructions document a one-time step that works on current macOS (removing the quarantine attribute, or System Settings → Privacy & Security → Open Anyway), and do not rely on the removed right-click → Open bypass

#### Scenario: Notarized build once the certificate exists
- **WHEN** a Developer-ID certificate and notary credentials are configured
- **THEN** the release produces a Developer-ID-signed, notarized, stapled build and the Gatekeeper caveat is dropped
