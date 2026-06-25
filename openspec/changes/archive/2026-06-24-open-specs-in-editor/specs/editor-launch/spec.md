## ADDED Requirements

### Requirement: Open the current spec in the external editor

The TUI SHALL allow the user to open the spec being viewed in the system editor by pressing `e` while in `ModeViewingSpec`. The editor SHALL be the value of the `$EDITOR` environment variable; if it is not defined, `vi` SHALL be used as a fallback. The file opened SHALL be `openspec/specs/<name>/spec.md` for the spec currently being viewed. This SHALL apply both when the full spec is rendered and when a single requirement is focused, since requirements are sections within the same `spec.md` file. The TUI SHALL suspend its control of the terminal before launching the editor using `tea.ExecProcess` and resume it on exit, with mouse tracking still functional after returning.

#### Scenario: Open full spec in editor

- **WHEN** the mode is `ModeViewingSpec` showing a full spec and the user presses `e`
- **THEN** the TUI yields the terminal and opens `$EDITOR` on that spec's `openspec/specs/<name>/spec.md`; when the editor is closed the TUI resumes with functional mouse tracking

#### Scenario: Open spec in editor while a requirement is focused

- **WHEN** the mode is `ModeViewingSpec` focused on a single requirement and the user presses `e`
- **THEN** the TUI opens `$EDITOR` on the same spec's `openspec/specs/<name>/spec.md` (the file containing that requirement)

#### Scenario: Fallback to vi when $EDITOR is not defined

- **WHEN** the mode is `ModeViewingSpec`, `$EDITOR` is not defined in the environment, and the user presses `e`
- **THEN** the TUI launches `vi` with the path of the spec being viewed

#### Scenario: Help bar advertises the edit shortcut in spec view

- **WHEN** the mode is `ModeViewingSpec`
- **THEN** the help bar includes `e: edit`

### Requirement: Reload spec content after editing

The TUI SHALL reload the content of the edited spec immediately upon returning from the editor, so that changes made externally are reflected without restarting the TUI. The view SHALL remain in `ModeViewingSpec`, preserving full-spec or requirement-focus state.

#### Scenario: Reload spec after editing

- **WHEN** the user edits a spec's `spec.md` in the editor and closes it while in `ModeViewingSpec`
- **THEN** the TUI re-renders the spec view with the updated content, staying in `ModeViewingSpec`

#### Scenario: Focused requirement reflects edits

- **WHEN** a single requirement is focused and the user edits that requirement's text in the editor and closes it
- **THEN** the focused requirement view shows the updated requirement content
