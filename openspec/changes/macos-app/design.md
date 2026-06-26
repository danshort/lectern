## Context

lectern splits cleanly into two layers: `internal/openspec` (pure domain — loader, `fileSystem` interface, task parsing, validation, git-worktree discovery; ~no UI coupling) and `internal/ui` (Bubble Tea TUI rendering ANSI to a terminal). The domain layer is reusable across any front end; the UI layer is terminal-only and reuses nowhere. A macOS app therefore reuses the *logic* (by re-implementation, since Swift can't import Go) and rebuilds the *presentation* from scratch.

The honest framing: there is **no risk to the existing TUI** (this change is purely additive). The domain port is a **bounded, mechanical** cost — but only partly pinnable by golden tests (see "What the corpus cannot pin"). The new app's **uncertainty and recurring cost** are concentrated in distribution, the App Sandbox/file-access model, and markdown/UX parity — *not* in the domain port. "No risk to the TUI" is not "no risk," and "biggest effort" (front end) is not the same as "biggest uncertainty" (sandbox + parity).

## Goals / Non-Goals

**Goals:**
- A native-feeling macOS reader (NavigationSplitView sidebar over changes → artifacts → rendered content) backed by a faithful port of the domain layer.
- A single source of truth for *expected behavior* shared by Go and Swift, enforced in CI — with its limits stated explicitly.
- Zero change to the existing TUI.

**Non-Goals:**
- cgo bridge; cross-platform GUI; artifact authoring beyond task toggling; shipping all phases at once.

## Decisions

- **Native Swift re-implementation, not a cgo bridge.** A `c-archive` bridge would *eliminate the entire drift surface* (one parser, no YAML/regex/sort/error-string divergence) and is cheaper than the original draft implied — the C boundary is "call one function, get the same JSON the golden already produces." It is rejected anyway, but on the *real* grounds: (a) `ToggleTask` writes files and `ListWorktrees` spawns `git` — putting the Go runtime in charge of process spawning and file writes inside a **sandboxed, signed, notarized** app fights App Sandbox and security-scoped bookmarks; (b) a statically-linked Go archive with its own runtime/threads complicates the hardened-runtime/entitlements story; (c) it forecloses native FSEvents. We accept the drift risk (contained by the golden corpus + the mitigations below) in exchange for a sandbox-friendly, FSEvents-native app. Revisit only if the OpenSpec format starts churning.
- **App Sandbox posture is an upstream architectural decision, resolved BEFORE Phase 3 (not in packaging).** A sandboxed app cannot freely `Process`-spawn `/usr/bin/git`, and `git worktree list` reads `.git` and sibling worktree dirs that are very likely **outside** the security-scoped bookmark granted by the folder picker. FSEvents and reveal-in-Finder/open-in-editor have the same scope constraint. We must choose: **Developer-ID, non-sandboxed** (worktrees + git work; no Mac App Store) vs **sandboxed** (App Store-eligible; worktree feature likely needs a helper or gets dropped). Defaulting to Developer-ID non-sandboxed unless App Store distribution is a goal. This choice dictates the entire file-access model, so it gates Phases 3–5.
- **Monorepo.** App lives under `macos/` so the corpus and goldens are one artifact both toolchains read. Cost: two toolchains; the Swift lane runs only on macOS runners.
- **`OpenSpecKit` as a UI-free SwiftPM target**, mirroring `internal/openspec` one-to-one — the same role `internal/openspec` plays for the TUI.
- **YAML is parsed with a real YAML library (Yams), not `Codable`.** `config.yaml` and `.openspec.yaml` are YAML; `Codable`/`JSONDecoder` does not parse YAML, and yaml.v3 edge cases (multiline scalars, `null` vs empty, coercion) must be matched. Decoding must also be **field-tolerant**: `.openspec.yaml` swallows all parse errors (→ empty `Created`), while `config.yaml` propagates unmarshal errors — an asymmetry a naive all-or-nothing `Codable` decode gets wrong.
- **FSEvents over polling.** The TUI polls disk at 500 ms; the app watches `openspec/` via `FSEventStream`/`DispatchSource`.
- **Markdown renderer: swift-markdown + custom SwiftUI views, not `AttributedString` alone.** OpenSpec artifacts use tables, fenced code blocks, and nested lists heavily; `AttributedString`'s markdown is single-paragraph-oriented and renders none of those well. Rendering also includes *behavior* the TUI performs (below), not just styling.
- **Phase 1 (corpus + Go golden test) ships first and stands alone**, hardening the existing loader regardless of whether the app is built.

### Swift module map (`OpenSpecKit` ⇆ `internal/openspec`)

| Go | Swift | Notes |
|---|---|---|
| `Artifact{Content,Present,ReadErr}` | `struct Artifact { content; present; readError? }` | `present` distinct from non-empty; `readError` is **not** byte-stably serializable (see corpus limits) |
| `NamedSpec` / `ProjectSpec` | `Codable` structs | carry placeholder `content` **and** a read-error flag, not content-only |
| `Change` / `Project` | `Codable` structs | |
| `ProjectConfig{Context, Rules}` | `struct { context; rules: [String:[String]]? }` | parsed via **Yams**; `nil` vs `[:]` (absent vs empty `rules:`) must match Go nil-vs-empty-map |
| `TaskItem{Kind,Text,Done,LineNum}` | `struct` + `enum ItemKind {section,task}` | `LineNum` indexes raw `\n`-split lines |
| `Worktree{...}` | `Codable struct` | |
| `fileSystem` interface | `protocol FileSystem` | injection seam; `readDir` **must sort by name** (see #13) |
| `Loader{fs}` | `struct Loader { let fs: FileSystem }` | default `OSFileSystem` |

### Foundational invariant (underpins most drift)

**Go `strings.Split` semantics.** `Split("a\nb\n","\n") == ["a","b",""]` (trailing empty for a trailing newline); `Split("","\n") == [""]` (length 1). This drives `splitLines` (parsing) *and* the toggle write path (split → mutate → `Join`). Swift's `components(separatedBy:)` matches; `split(separator:)` does **not** (drops trailing/empty) — and a porter reaches for `split` first. Getting this wrong shifts every `LineNum`, makes `ToggleTask` index the wrong line, and adds/drops a trailing newline → byte divergence on every write. Port to `components(separatedBy:)` and fixture a trailing-newline file.

### Behaviors that are the actual drift risk

Verified against the Go source; the corpus must pin each (corrections from review folded in):

1. **Change sort** (`loader.go:129`): `Created` descending, empty `Created` last, **stable**. The `Name`-ascending tiebreak applies **only when BOTH `Created` are empty**; two equal *non-empty* `Created` keep input order via stability (no name tiebreak). Use an index tiebreak; do **not** add a name tiebreak to the non-empty-equal case.
2. **CRLF asymmetry (the sharp edge):** `splitLines` strips trailing `\r`, but `ToggleTask` splits raw `"\n"` and rejoins on `"\n"` — CRLF files stay CRLF on write. Port the two paths separately; byte-exact LF *and* CRLF toggle goldens.
3. **`ToggleTask` is line-number indexed**, replacing the *first* `"- [x] "`/`"- [ ] "` substring on `items[idx].LineNum`; bounds checks (`idx>=len || lineNum>=len → no-op`) port exactly. It also **mutates the caller's `items[idx].Done` in place** (Go slice aliasing) — a Swift value-type `[TaskItem]` loses this unless `inout`; the in-memory flip is relied on by callers and is invisible to both goldens.
4. **Artifact error semantics:** not-found → absent; other read error → `Present:true`, `Content = "⚠ couldn't read " + path + ": " + err.Error()`. `ValidateChange` skips specs with a read error.
5. **Archive name parse:** regex `^(\d{4}-\d{2}-\d{2})-(.+)$` *and* the prefix must parse as a real calendar date (`2006-01-02`), else `(dir, "")`. So `2026-13-99-foo` and `2026-02-29-foo` are rejected; `2024-02-29-foo` is accepted — match Go's calendar validity, not just the regex. Archive list reverse-sorted (not stable; fine, names unique).
6. **Spec aggregation:** join `"# "+name+"\n\n"+content` with `"\n\n---\n\n"`; empty → absent.
7. **`ExtractRequirement`:** block from the matching `### Requirement:` (exact trimmed name) to the next `### Requirement:`.
8. **Validation (expanded):** `ValidateChange` also requires `proposal.Present` (`"missing proposal.md"`) — not just delta rules. Header detection uses **`HasPrefix`**, so `## Purpose` matches `## Purpose and Scope` and any `## ` line closes a requirement's scenario-search window; but `deltaHeaderRe` is **anchored** (`^## (ADDED|…) Requirements\s*$`, exact, singular `Requirement` fails). Do not unify the two matching disciplines. Empty-named requirements (`### Requirement:` with nothing after) are silently not scenario-validated.
9. **Worktree:** parse `git worktree list --porcelain` (5 s timeout via a `Process` watchdog); `markCurrentWorktree` via `rev-parse --show-toplevel`. `normalizePath` = `filepath.EvalSymlinks` **with a fall-back to lexical `filepath.Clean` when the path does not exist** — `URL.resolvingSymlinksInPath` does *not* replicate this (it doesn't error on nonexistent paths and normalizes `/private` differently). Port the EvalSymlinks-then-Clean-fallback explicitly; test a nonexistent path.

Other non-obvious behaviors to port and fixture: **unsorted `loadSpecs`/`ReloadChange` ordering** — these rely on Go `os.ReadDir`'s filename-sort guarantee, which Swift `FileManager.contentsOfDirectory` does **not** provide, so `FileSystem.readDir` must sort (#13, highest-value after the Split invariant; only visible with ≥2 spec dirs); **two spec-loading paths with different error semantics** (`loadSpecs` swallows `listDirs` errors → empty; `LoadProjectSpecsFrom` propagates them); **`LoadFromPath`** derives `Project.Name` from the *grandparent* dir and requires `.openspec.yaml`; the **`⚠` placeholder is multibyte** (`U+26A0`) so any byte-offset math on `Content` diverges between Go (bytes) and Swift (graphemes).

### Shared golden harness

```
testdata/corpus/
  basic-project/openspec/...        # ≥3 changes (mixed/equal Created) + ≥3 spec dirs — exercises sort stability & loadSpecs order
  crlf-tasks/ , lf-tasks/           # both toggle write goldens (CRLF + LF), trailing-newline included
  unreadable-artifact/...           # error-normalized golden (see below)
  malformed-archive-name/...        # 2026-13-99-foo, 2026-02-29-foo, 2024-02-29-foo (calendar validity)
  malformed-meta/...                # bad .openspec.yaml → empty Created (tolerant decode)
  config-variants/...               # absent rules, rules:{}, multiline context (nil vs empty, YAML semantics)
  delta-specs/...                   # validation incl. missing proposal, HasPrefix headers, empty-named req
  worktree-porcelain/*.txt          # captured porcelain text → parser golden (no live git)
  golden/
    *.json            # serialized Project, sorted keys
    tasks.json        # ParseTasks output incl. LineNum
    requirements.json # ExtractRequirement
    worktrees.json    # parseWorktreeList over the captured text
    config.md         # ConfigToMarkdown
    validation.json   # path → []messages
    *.after-toggle.tasks.md  # byte-exact write goldens (LF + CRLF)
```

- **Go side:** a golden test runs each entry point over the corpus, serializes with sorted keys, compares to `golden/*` (with `-update`). Adds to — does not replace — the existing `t.TempDir()` unit tests.
- **Swift side:** `OpenSpecKitTests` runs the same entry points against the same corpus, encodes with `.sortedKeys`, asserts byte-equality.
- **Explicit encoding contract** (where `Codable` and `encoding/json` most often disagree): no `omitempty`; encode `nil`/absent as `null`; encode empty slices as `[]` (never `null`); empty maps as `{}`; field names snake_case via Go tags ↔ Swift `CodingKeys`. Without this written contract, nil-vs-empty fields produce false diffs even when logic is identical.
- **Error-string normalization:** the unreadable-artifact placeholder embeds `err.Error()`, which is OS/locale-specific and differs between Go and a Swift `NSError`. Its golden stores `present:true` + a read-error flag + a **prefix-only** content match (`"⚠ couldn't read <path>: "`), not the raw error tail. Otherwise this golden is permanently red or skipped — a real defect in the naive byte-exact design.
- **The harness covers** Project, tasks, requirements, worktree-parsing, config, validation, and the toggle write path. It does **not** cover ToggleTask's in-memory `items` mutation (#3) → assert that in unit tests on both sides.

### What the corpus cannot pin (the residual drift surface)

Golden coverage ≠ behavioral completeness. These have no stable golden and need targeted unit/integration tests + manual QA:
- **Cross-platform FS/Unicode:** macOS APFS filename normalization (NFD) vs Linux CI; case-insensitive-preserving filesystems; `normalizePath` symlink/`/private` behavior is **macOS-only** and the Go golden runs on Linux — it cannot be exercised there at all.
- **YAML semantics:** yaml.v3 vs Yams on multiline scalars, anchors, coercion, `null`.
- **Regex dialect:** Go RE2 vs Swift ICU (`NSRegularExpression`/`Regex`) on `\s`, anchoring, Unicode classes — corpus pins known inputs, not the dialect gap.
- **OS error-string formatting** (handled by the normalization above).
- **FSEvents timing, SwiftUI views, live `git` invocation/timeout** (the porcelain *parser* is golden-able; the `Process` plumbing is not).

A rule: **any change to `internal/openspec` loader/tasks/validate/worktree-parser behavior must add or modify a fixture.** Enforced by checklist/CI nudge; the corpus is the canonical behavior spec.

### Markdown/UX behaviors the app must replicate (not just styling)

The TUI's `viewport.go` does more than render: it **injects a validation banner** (`> ⚠ **Validation errors**`) into spec markdown before rendering (skipped for unreadable specs); it does **per-requirement focus** via `ExtractRequirement` and **jump-to-line** by substring-matching rendered output. These are requirements for the app, not parity niceties — captured in `macos-app/spec.md`.

## Risks / Trade-offs

- **[High] Effort and uncertainty are front-end/distribution-dominated.** The reusable domain logic is the small, now well-understood part. SwiftUI, FSEvents, git plumbing, sandbox/entitlements, and packaging are net-new and carry the uncertainty.
- **[High] App Sandbox vs `git`/arbitrary-path access** (see Decisions) — resolve before Phase 3 or risk reworking the file-access model at Phase 5.
- **[Medium] Logic drift Go↔Swift**, partly unpinnable (FS/Unicode, YAML, regex dialect). Mitigated by corpus + the fixture-on-change rule; residual risk acknowledged.
- **[Medium] Signing + notarization** require an Apple Developer Program account (**$99/yr recurring**), cert rotation, and notary credentials in CI; notarization is an Apple network service that can fail/slow a release.
- **[Medium] Markdown renderer** is real work (swift-markdown + custom views) plus a dependency/licensing decision; visual parity is approximate and not golden-tested.
- **[Medium] Ongoing two-toolchain tax:** every OpenSpec format change becomes a 2-implementation + corpus change; contributors must keep Swift green. Additive at runtime, not at maintenance.
- **[Medium] Accessibility** (VoiceOver, Dynamic Type, keyboard nav, contrast) and **app lifecycle** (auto-update via Sparkle vs `brew upgrade`, recent-projects, settings, state restoration) are baseline Mac expectations, currently unscoped.
- **[Medium] Toggle race in a live-reload GUI:** `items[idx].LineNum` captured at render can be stale if the file changed on disk; re-read + re-parse immediately before toggling (or hash-check).
- **[Low] CI cost:** the macOS Swift golden lane must run **unconditionally** (no path filter) or the anti-drift guarantee is hollow; pin the Xcode/Swift toolchain for byte-stable goldens. macOS runner minutes ~10× Linux.
- **[Low] Release coupling:** the macOS sign/notarize job must be **decoupled** from the existing goreleaser/CLI publish so a notary outage never blocks `brew install lectern`.

## Migration / Rollout

Phased, each independently valuable; abandon after any phase with nothing lost. **Phases 3–5 are gated on the App Sandbox + markdown-renderer decisions.**
1. **Corpus + Go golden test** — hardens the existing loader; no Swift yet.
2. **`OpenSpecKit`** — port + Swift tests green against the shared golden (Yams, encoding contract, error normalization).
3. **SwiftUI shell** — read-only browse + native markdown (incl. validation banner, requirement focus/jump) over `OpenSpecKit`. *Gated on sandbox + renderer decisions.*
4. **Interaction + OS integration** — CRLF-safe task toggle (re-read before write), worktrees view, FSEvents live reload.
5. **Packaging** — sign, notarize, decoupled cask/release lane, accessibility pass.
