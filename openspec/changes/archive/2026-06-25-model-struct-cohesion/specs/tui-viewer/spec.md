## ADDED Requirements

### Requirement: Mode transitions enforce state invariants centrally

The TUI SHALL change its active mode through a single transition function rather than assigning the mode field directly at each call site. That function SHALL be the only place the mode is assigned, and on every transition it SHALL enforce the state invariants the destination mode owns: it SHALL clamp the destination mode's cursors into the valid range for the current data, and it SHALL reset state that belongs only to the mode being left so it cannot leak into the next mode. In particular, leaving the spec-viewing mode SHALL clear the requirement focus state (focus flag and jump target). The transition function SHALL NOT invalidate the rendered-artifact cache; cache invalidation remains the responsibility of the individual call sites, which keep the cache across tab switches by design.

This requirement constrains internal state management only. It introduces no change to user-visible navigation, key bindings, or rendered output: every navigation behavior specified elsewhere in this capability continues to hold unchanged.

#### Scenario: Leaving spec focus view clears focus state

- **WHEN** the user is viewing a requirement in focus mode (spec-viewing mode with a focus flag and jump target set) and transitions to any other mode
- **THEN** the focus flag and jump target are cleared, so a later entry into the spec-viewing mode that is not a focused requirement does not inherit the previous jump target

#### Scenario: Entering a tabbed mode clamps to an available tab

- **WHEN** a transition enters a mode that renders artifact tabs (normal or archive viewing) while the previously selected tab is not available on the destination change
- **THEN** the selected tab is clamped to the first available tab, and no transition leaves the selected tab pointing at an absent artifact

#### Scenario: Entering the spec-viewing mode clamps the spec cursor

- **WHEN** a transition enters the spec-viewing mode with a spec cursor that is out of range for the loaded project specs
- **THEN** the spec cursor is clamped into range, and an empty project-spec list resolves to a non-negative cursor with no out-of-range access

#### Scenario: Transitions do not drop the render cache

- **WHEN** the user switches between artifact tabs and returns to a previously rendered tab
- **THEN** the previously rendered content is served from cache without re-rendering, because the mode-transition function does not clear the render cache
