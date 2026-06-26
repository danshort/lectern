## Context

lectern splits cleanly into two layers: `internal/openspec` (pure domain — loader, `fileSystem` interface, task parsing, validation, git-worktree discovery; ~no UI coupling) and `internal/ui` (Bubble Tea TUI rendering ANSI to a terminal). The domain layer is reusable across any front end; the UI layer is terminal-only and reuses nowhere. A macOS app therefore reuses the *logic* (by re-implementation, since Swift can't import Go) and rebuilds the *presentation* from scratch.

The honest framing: there is **no functionality risk** to the TUI (this change is purely additive), the domain-port carries a **bounded, mechanical** cost (pinned by golden tests), and the SwiftUI + OS-integration layer is **open-ended net-new work** with no golden to lean on. Total effort is dominated by the part that does *not* reuse.

## Goals / Non-Goals

**Goals:**
- A native-feeling macOS reader (NavigationSplitView sidebar over changes → artifacts → rendered content) backed by a faithful port of the domain layer.
- A single source of truth for *expected behavior* shared by Go and Swift, enforced in CI.
- Zero change to the existing TUI.

**Non-Goals:**
- cgo bridge; cross-platform GUI; artifact authoring; shipping all phases at once.

## Decisions

- **Native Swift re-implementation, not a cgo bridge.** Bridging `internal/openspec` via `go build -buildmode=c-archive` keeps one source of truth but drags the Go runtime + per-arch fat static libs into an otherwise pristine Mac app and forces struct marshaling across a C boundary. The domain surface is small and stable (a handful of conventions + a task regex), so a native `Codable` port is idiomatic, testable, and lighter. Drift is contained by the shared golden corpus instead. Revisit only if the OpenSpec format starts churning.
- **Monorepo.** App lives under `macos/` in this repo so the fixture corpus and golden files are a single artifact both toolchains read, with one CI. Cost: two toolchains in one repo.
- **`OpenSpecKit` as a UI-free SwiftPM target.** Mirrors `internal/openspec` one-to-one and is the only thing the app and the Swift tests depend on — the same role `internal/openspec` plays for the TUI.
- **FSEvents over polling.** The TUI polls disk at 500 ms; the app should watch `openspec/` via `FSEventStream`/`DispatchSource`. A place where native legitimately beats the port.
- **Phase 1 (corpus + Go golden test) ships first and stands alone.** It hardens the existing loader regardless of whether the app is ever built, so it is low-regret.

### Swift module map (`OpenSpecKit` ⇆ `internal/openspec`)

| Go | Swift | Notes |
|---|---|---|
| `Artifact{Content,Present,ReadErr}` | `struct Artifact { content; present; readError? }` | keep `present` distinct from non-empty |
| `NamedSpec` / `ProjectSpec` / `Change` / `Project` | `Codable` structs | Codable is what makes golden tests trivial |
| `ProjectConfig{Context, Rules}` | `struct { context; rules: [String:[String]] }` | |
| `TaskItem{Kind,Text,Done,LineNum}` | `struct` + `enum ItemKind {section,task}` | |
| `Worktree{...}` | `Codable struct` | |
| `fileSystem` interface | `protocol FileSystem` | injection seam for tests |
| `Loader{fs}` | `struct Loader { let fs: FileSystem }` | default `OSFileSystem` |

### Behaviors that are the actual drift risk

These are the non-obvious points where a naive port silently differs; the corpus must pin each:

1. **Change sort** (`LoadFrom`): `Created` descending, empty `Created` last, ties by `Name` ascending, **stable**. Swift `sort` isn't stable — use an index tiebreak.
2. **CRLF asymmetry (the sharp edge):** `splitLines` strips trailing `\r` for parsing/matching, but `ToggleTask` splits raw `"\n"` *without* stripping `\r`, so CRLF files stay CRLF on write. Porting toggle with the stripping helper would rewrite every line ending. Port the two paths separately; test a CRLF fixture with a byte-exact golden.
3. **`ToggleTask` is line-number indexed**, replacing the first `"- [x] "`/`"- [ ] "` on `items[idx].LineNum`; bounds checks (`idx>=len || lineNum>=len → no-op`) port exactly.
4. **Artifact error semantics:** not-found → absent (`Artifact{}`); any other read error → `Present:true` with `Content = "⚠ couldn't read " + path + ": " + err`. `ValidateChange` skips specs with a read error (read failure ≠ structural failure).
5. **Archive name parse:** regex `^(\d{4}-\d{2}-\d{2})-(.+)$` *and* the prefix must parse as a real `2006-01-02` date, else `(dir, "")`. Archive list sorted reverse.
6. **Spec aggregation:** join `"# "+name+"\n\n"+content` with separator `"\n\n---\n\n"`; empty → absent.
7. **`ExtractRequirement`:** block from the matching `### Requirement:` (exact trimmed name) to the next `### Requirement:`.
8. **Validation:** spec needs `## Purpose` + `## Requirements` + every `### Requirement:` has a `#### Scenario:`; delta spec needs an `## ADDED|MODIFIED|REMOVED|RENAMED Requirements` header, scenarios required only under ADDED/MODIFIED.
9. **Worktree:** parse `git worktree list --porcelain` (5s timeout via `Process` watchdog); `markCurrentWorktree` via `rev-parse --show-toplevel`; `normalizePath` resolves symlinks (`/var`↔`/private/var`) — use `URL.resolvingSymlinksInPath`.

### Shared golden harness

```
testdata/corpus/
  basic-project/openspec/...        # changes, specs, config.yaml
  crlf-tasks/openspec/...           # the \r torture case
  unreadable-artifact/...           # placeholder/⚠ path
  malformed-archive-name/...        # 2026-13-99-foo
  delta-specs/...                   # validation cases
  golden/
    basic-project.json              # serialized Project (sorted keys)
    crlf-tasks.after-toggle.tasks.md# byte-exact write output
    validation.json                 # path → []messages
```

- **Go side:** a golden test walks the corpus, runs the loader, `json.Marshal`s with sorted keys, compares to `golden/*.json` (with `-update` to regenerate). Adds to — does not replace — the existing `t.TempDir()` unit tests.
- **Swift side:** `OpenSpecKitTests` points `Loader` at the same corpus, encodes with `JSONEncoder(.sortedKeys)`, asserts byte-equality against the same golden files. Field names aligned via Go struct tags ↔ Swift `CodingKeys` (snake_case).
- **Contract holds** because both fixtures and expected output are identical bytes; CI runs both lanes, so a one-sided rule change is a red build, not a latent bug.
- **The toggle write path** is golden'd as post-toggle *file bytes* (not JSON) — the only way to catch the CRLF regression.

### What has no golden (the real cost)

Parsing reduces to bytes-in/value-out, so it's cheap to pin. These do not:
- **FSEvents watcher** — correctness is "UI refreshed at the right time," a timing behavior, not a value.
- **SwiftUI views** — output is pixels/interaction.
- **Live `git` invocation + timeout** — the porcelain *parser* is golden-able from captured text, but the `Process` plumbing and 5s watchdog need integration tests.

These fall to integration tests + manual QA — slower, less certain, and the bulk of a from-scratch native app. A cgo bridge would **not** save this layer; it only shares domain logic.

## Risks / Trade-offs

- **[High] Effort is front-end-dominated.** The reusable part (domain logic) is the small part. The SwiftUI app, FSEvents, git plumbing, and packaging are net-new and where most time goes. This plan's value is making that explicit before committing.
- **[Medium] Logic drift between Go and Swift.** Mitigated by the shared corpus + golden CI on both lanes; the residual risk is behavior with no golden (watcher/UI).
- **[Medium] Signing + notarization.** Requires an Apple Developer account and signing secrets in CI; the one genuinely new operational cost. Without it, Gatekeeper blocks a smooth install.
- **[Low] Monorepo weight.** Two toolchains (Go + Swift/Xcode) in one repo; Swift build only runs on macOS runners. Contained to the `macos/` lane.
- **[Low] Markdown rendering parity.** TUI uses Glamour (markdown→ANSI); the app renders via `AttributedString`/a Swift markdown lib — visual parity is approximate and not golden-tested.

## Migration / Rollout

Phased, each independently valuable; abandon after any phase with nothing lost:
1. **Corpus + Go golden test** — hardens the existing loader; no Swift yet.
2. **`OpenSpecKit`** — port + Swift tests green against the shared golden.
3. **SwiftUI shell** — read-only browse + native markdown over `OpenSpecKit`.
4. **Interaction + OS integration** — task toggle (CRLF-safe), worktrees view, FSEvents live reload.
5. **Packaging** — sign, notarize, cask, CI lane.
