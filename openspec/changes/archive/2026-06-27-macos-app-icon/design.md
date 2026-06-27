## Context

The app is a manually-assembled SwiftPM bundle (no Xcode project), so the icon
is compiled in `package.sh` rather than by an Xcode build.

## Decisions

- **`actool` compiles the `.icon` directly** — `xcrun actool lectern.icon
  --compile <Resources> --app-icon lectern …` emits `Assets.car` (the Liquid
  Glass icon) + `lectern.icns` (fallback) and the `CFBundleIconName/File` keys.
  Verified with Xcode 26.6.
- **Compile before signing** so the icon is covered by the signature.
- **Graceful on old toolchains** — Xcode < 26's `actool` can't read `.icon`;
  the step warns and continues so a release on an older runner still builds
  (icon-less). The release workflow prefers the newest installed Xcode.

## Risks / Trade-offs

- **[Low] CI runner Xcode** — released builds only get the icon when the runner
  has Xcode 26+; otherwise icon-less but functional. Local builds (Xcode 26) are
  fine.
