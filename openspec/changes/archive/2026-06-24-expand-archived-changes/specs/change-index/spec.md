## ADDED Requirements

### Requirement: Expand an archived change with Space
While the mode is `ModeIndex` and the cursor is on an archived change, pressing `Space` SHALL toggle the expansion of that archived change. When expanded, the index SHALL insert, immediately after the archived change row, one nested sub-item for each artifact type that is present on that change, in the fixed order proposal, design, specs, tasks. Artifact types that are not present SHALL be omitted. Pressing `Space` again SHALL collapse the archived change and remove its sub-items. After toggling, the cursor SHALL remain anchored on the toggled archived change.

#### Scenario: Expanding an archived change reveals its present artifacts
- **WHEN** the mode is `ModeIndex`, the cursor is on an archived change that has a proposal, specs, and tasks (but no design), and the user presses `Space`
- **THEN** three nested sub-items labelled "proposal", "specs", and "tasks" appear immediately below the archived change, in that order, and "design" is not shown

#### Scenario: Collapsing an archived change hides its artifacts
- **WHEN** the mode is `ModeIndex`, an archived change is expanded, and the user presses `Space` on it again
- **THEN** the nested artifact sub-items are removed and only the archived change row remains

#### Scenario: Cursor stays on the toggled archived change
- **WHEN** the mode is `ModeIndex`, the cursor is on an archived change, and the user presses `Space`
- **THEN** the cursor remains on that same archived change row

#### Scenario: Space on an archived change with no artifacts does nothing visible
- **WHEN** the mode is `ModeIndex`, the cursor is on an archived change that has no present artifacts, and the user presses `Space`
- **THEN** no sub-items are added and the cursor does not move

### Requirement: Navigate and display archived artifact sub-items
Archived artifact sub-items SHALL be navigable with `j` (down) and `k` (up) like any other index item, and SHALL be rendered indented below their parent archived change with the artifact-type name. The sub-item under the cursor SHALL be visually highlighted using the same cursor style as requirement sub-items.

#### Scenario: Navigate from archived change into its artifacts
- **WHEN** the mode is `ModeIndex`, an archived change is expanded, the cursor is on that archived change, and the user presses `j`
- **THEN** the cursor moves to the first artifact sub-item below it

#### Scenario: Artifact sub-item under the cursor is highlighted
- **WHEN** the mode is `ModeIndex` and the cursor is on an archived artifact sub-item
- **THEN** that sub-item is rendered indented with the highlighted cursor marker

### Requirement: Open an archived artifact with Enter
Pressing `Enter` on an archived artifact sub-item, or left-clicking it when already selected, SHALL switch the mode to `ModeViewingArchive` for the parent archived change with the active tab set to the selected artifact type.

#### Scenario: Enter on an artifact sub-item opens that tab
- **WHEN** the mode is `ModeIndex`, the cursor is on the "design" sub-item of an expanded archived change, and the user presses `Enter`
- **THEN** the mode switches to `ModeViewingArchive` for that archived change with the active tab set to `design`

#### Scenario: Click on a selected artifact sub-item opens that tab
- **WHEN** the mode is `ModeIndex` and the user left-clicks the already-selected "tasks" sub-item of an expanded archived change
- **THEN** the mode switches to `ModeViewingArchive` for that archived change with the active tab set to `tasks`

### Requirement: Filtering keeps archived artifact sub-items with matching parents
When a filter is active, an archived artifact sub-item SHALL be considered a match when its parent archived change name matches the filter or its artifact-type label matches the filter, so that the sub-items of a matching archived change remain visible while that change is expanded.

#### Scenario: Sub-items remain visible when the parent change matches
- **WHEN** the mode is `ModeIndex`, an archived change named "data-export" is expanded, and the user types `/data`
- **THEN** the "data-export" archived change and its visible artifact sub-items remain shown
