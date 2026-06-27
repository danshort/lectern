## 1. Icon in the bundle

- [x] 1.1 Add the Icon Composer source at `macos/LecternApp/lectern.icon`
- [x] 1.2 `Info.plist`: `CFBundleIconName`/`CFBundleIconFile` = `lectern`
- [x] 1.3 `package.sh`: compile the `.icon` via `actool` into `Contents/Resources` before signing; warn + continue if the toolchain can't (Xcode < 26)
- [x] 1.4 `macos-app-release.yml`: select the newest installed Xcode (26+ for `.icon`)

## 2. Spec + verification

- [x] 2.1 Delta spec: ADD an "App icon" requirement to `macos-app`
- [x] 2.2 `package.sh 0.1.0` produces a bundle containing `Assets.car` + `lectern.icns` with `CFBundleIconName` set; app still launches
- [x] 2.3 Manual QA: the packaged app shows the icon in Finder / Dock
