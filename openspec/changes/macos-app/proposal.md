## Why

lectern is a terminal-only reader for OpenSpec artifacts. A native macOS app would give the same "window into your specs" to people who'd rather click than live in a TUI, and would bundle naturally with the CLI (`brew install lectern` → CLI + app). The opportunity is that the correctness-critical part of lectern — the `internal/openspec` domain layer (loader, task parsing, validation, worktree discovery) — is already UI-agnostic, so a second front end is additive and carries **no risk to the existing TUI**.

This change is a feasibility-grade plan, not a commitment to ship. It exists to make the cost honest: what reuses, what must be rebuilt, and where the real effort lands.

## What Changes

- Add a native **SwiftUI macOS app** in this repo (monorepo) under `macos/`, reading the same `openspec/` layout the TUI reads.
- Add **`OpenSpecKit`**, a UI-free Swift package that re-implements the `internal/openspec` domain layer (parsing, tasks, validation, worktree porcelain). Native Swift, not a cgo bridge — see `design.md` for why.
- Add a **shared fixture corpus** (`testdata/corpus/`) plus committed **golden output**, asserted by *both* a new Go golden test and the Swift test suite, so the two implementations cannot silently drift.
- Extend release/CI to build, sign, and notarize the `.app`/`.dmg` and publish it (Homebrew cask) alongside the existing CLI artifacts.
- The Go TUI is **untouched**: no refactor, no extraction, no behavior change.

## Non-goals

- Rewriting or replacing the TUI. The CLI remains the primary tool.
- A cgo / `c-archive` bridge to share Go logic into Swift (rejected in `design.md`; native re-implementation chosen instead).
- Editing artifacts beyond what the TUI already does (toggle task checkboxes). No proposal/spec authoring in the app.
- Cross-platform GUI (Windows/Linux). This is macOS-native by intent.
- Shipping in one release. Phases land independently; the corpus + golden harness (Phase 1) is useful on its own even if the app never ships.

## Capabilities

### New Capabilities

- `shared-fixture-corpus`: a committed corpus of OpenSpec project fixtures with golden output that every implementation of the loader (Go and Swift) must reproduce, run in CI for both.
- `macos-app`: a native macOS reader for OpenSpec artifacts mirroring the TUI's read/navigate/toggle/worktree/live-reload behaviors.

## Impact

- `testdata/corpus/**` — new shared fixtures + `golden/` output (also hardens the existing Go loader).
- `internal/openspec/golden_test.go` — new Go test asserting the loader reproduces the golden output.
- `macos/OpenSpecKit/**` — new Swift package (domain layer port + tests).
- `macos/LecternApp/**` — new SwiftUI app target.
- `.github/workflows/**` — new macOS build/sign/notarize lane; cask publication.
- `.goreleaser.yaml` — extended (or paired tooling) to attach the `.app`/`.dmg`.
- One-time operational cost: an Apple Developer account + signing secrets in CI (notarization).
- TUI / `internal/ui` / `cmd/lectern`: **no change**.
