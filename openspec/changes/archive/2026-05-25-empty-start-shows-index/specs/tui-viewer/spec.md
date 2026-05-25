## MODIFIED Requirements

### Requirement: Pantalla de bienvenida sin changes activos
When the TUI starts and there are no active changes, it SHALL open directly in the index view (`ModeIndex`), showing the "Active Changes" section (empty), "Specifications", and "Archived Changes" sections within the full TUI chrome, with the index help bar. If the TUI enters `ModeNormal` while there are no active changes (e.g., all changes were deleted during the session), it SHALL show an informational message with the available actions.

#### Scenario: Arranque sin changes activos muestra el índice
- **WHEN** the TUI starts and `openspec/changes/` contains no active subdirectories
- **THEN** the TUI shows the index view with "Active Changes", "Specifications", and "Archived Changes" sections and the help bar `j/k: navigate  Enter: open  Space: expand  s: sort by suffix  i: info  Esc: quit`

#### Scenario: Sin changes activos desde ModeNormal
- **WHEN** the mode is `ModeNormal` and there are no active changes
- **THEN** the TUI shows `"No active changes. Create one with /opsx:propose"` and the help line `a/Esc: index  q: quit`
