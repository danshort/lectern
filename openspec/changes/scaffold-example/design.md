## Context

This change exists to exercise the macOS app's multi-change UI. It mirrors the
shape of a real change (proposal, design, tasks, a delta spec) without touching
any code.

## Decisions

- **Keep it minimal but complete.** One capability, one requirement with one
  scenario, a couple of tasks — enough to render every artifact type.

## Risks / Trade-offs

- **[Low] Demo clutter.** It adds a second active change to the repo; remove it
  before the integration branch merges to `main`.
