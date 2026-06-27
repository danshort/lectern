#!/usr/bin/env bash
# Regenerates the committed app-icon assets from the Icon Composer source.
#
# The `.icon` format needs a recent actool (Xcode 26.4+). CI runners lag, so we
# commit the compiled output (Resources/AppIcon/{Assets.car,lectern.icns}) and
# package.sh ships those. Run this whenever lectern.icon changes, then commit
# the regenerated Resources/AppIcon.
set -euo pipefail

PKG_DIR="$(cd "$(dirname "$0")/.." && pwd)"   # macos/LecternApp
ICON="$PKG_DIR/lectern.icon"
OUT="$PKG_DIR/Resources/AppIcon"

[ -d "$ICON" ] || { echo "missing $ICON" >&2; exit 1; }

tmp="$(mktemp -d)"
if ! xcrun actool "$ICON" \
        --compile "$tmp" \
        --app-icon lectern \
        --platform macosx \
        --minimum-deployment-target 13.0 \
        --output-partial-info-plist "$tmp/icon-partial.plist" >/dev/null 2>&1 \
        || [ ! -f "$tmp/lectern.icns" ]; then
    echo "ERROR: actool could not compile $ICON — needs Xcode 26.4+ (you have:" >&2
    echo "       $(xcodebuild -version 2>/dev/null | head -1))." >&2
    exit 1
fi

mkdir -p "$OUT"
cp "$tmp/Assets.car" "$tmp/lectern.icns" "$OUT/"
echo "Regenerated $OUT/{Assets.car,lectern.icns} — commit these."
