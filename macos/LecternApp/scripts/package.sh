#!/usr/bin/env bash
# Packages LecternApp as a .app bundle, ad-hoc code-signs it (required for Apple
# Silicon to run an unsigned binary), and zips it for the Homebrew cask.
# Developer-ID signing + notarization are deferred (no Apple Developer account).
#
# Usage: scripts/package.sh [version]   (default version: 0.0.0-dev)
set -euo pipefail

VERSION="${1:-0.0.0-dev}"
APP_NAME="Lectern"
EXEC_NAME="LecternApp"
PKG_DIR="$(cd "$(dirname "$0")/.." && pwd)"   # macos/LecternApp
DIST="$PKG_DIR/dist"
APP="$DIST/$APP_NAME.app"

echo "==> swift build -c release"
swift build -c release --package-path "$PKG_DIR"
BIN="$(swift build -c release --package-path "$PKG_DIR" --show-bin-path)/$EXEC_NAME"

echo "==> assembling $APP_NAME.app"
rm -rf "$APP"
mkdir -p "$APP/Contents/MacOS" "$APP/Contents/Resources"
cp "$BIN" "$APP/Contents/MacOS/$EXEC_NAME"
sed "s/__VERSION__/$VERSION/g" "$PKG_DIR/Resources/Info.plist" > "$APP/Contents/Info.plist"

echo "==> ad-hoc code-signing"
codesign --force --deep --options runtime --sign - "$APP"
codesign --verify --strict --verbose=2 "$APP"

echo "==> zipping for cask"
ZIP="$DIST/$APP_NAME-$VERSION.zip"
rm -f "$ZIP"
( cd "$DIST" && ditto -c -k --keepParent "$APP_NAME.app" "$(basename "$ZIP")" )

echo ""
echo "Built:  $APP"
echo "Zip:    $ZIP"
echo "SHA256: $(shasum -a 256 "$ZIP" | cut -d' ' -f1)"
