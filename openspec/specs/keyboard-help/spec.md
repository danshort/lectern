# keyboard-help Specification

## Purpose
Defines a `?`-triggered keyboard-shortcut help overlay for the TUI: a centered, modal-like box that lists shortcuts grouped by screen, can be opened from any mode (except an active index filter input), and is dismissed to restore the originating screen and state without quitting.

## Requirements
### Requirement: Open keyboard-help overlay with `?`
The TUI SHALL open a keyboard-shortcut help overlay when the user presses `?`. The `?` key SHALL be honored from every mode: `ModeNormal`, `ModeIndex`, `ModeViewingArchive`, `ModeViewingSpec`, and `ModeViewingConfig`. Opening the overlay SHALL NOT change the underlying mode or any view state; it records the originating mode so it can be restored on dismissal.

#### Scenario: Open from the change viewer
- **WHEN** the mode is `ModeNormal` and the user presses `?`
- **THEN** the keyboard-help overlay is shown and the underlying mode remains `ModeNormal`

#### Scenario: Open from the index
- **WHEN** the mode is `ModeIndex` (and the filter input is not active) and the user presses `?`
- **THEN** the keyboard-help overlay is shown and the underlying mode remains `ModeIndex`

#### Scenario: Open from a spec, archive, or config view
- **WHEN** the mode is `ModeViewingSpec`, `ModeViewingArchive`, or `ModeViewingConfig` and the user presses `?`
- **THEN** the keyboard-help overlay is shown and the underlying mode is unchanged

### Requirement: Overlay lists shortcuts grouped by screen
While the overlay is open, the TUI SHALL render a centered, bordered, modal-like box over the current screen. The box SHALL list keyboard shortcuts as key/description pairs organized into per-screen groups, each with a visible group heading. The groups SHALL cover the global shortcuts and each screen: Global, Index, Change viewer, Archive viewer, Spec viewer, and Config viewer. Every shortcut documented SHALL match the binding actually handled by that screen.

#### Scenario: Overlay shows screen groups
- **WHEN** the keyboard-help overlay is rendered
- **THEN** it contains a distinct group heading for the Global shortcuts and for each screen (Index, Change viewer, Archive viewer, Spec viewer, Config viewer)

#### Scenario: Overlay documents representative shortcuts
- **WHEN** the keyboard-help overlay is rendered
- **THEN** the Index group lists `j`/`k` (navigate), `Enter` (open), and `/` (filter); the Change viewer group lists `1`-`4` (artifact tabs), `h`/`l` (change), and `Space` (toggle task); and the Global group lists `?` (help) and `q` (quit)

#### Scenario: Overlay is centered over the current view
- **WHEN** the terminal size is known and the overlay is rendered
- **THEN** the box is horizontally and vertically centered within the terminal width and height

### Requirement: Dismiss overlay restores the originating screen
The TUI SHALL close the keyboard-help overlay when the user presses `?`, `Esc`, or `q` while it is open, returning to the screen and state from which it was opened. Dismissing the overlay SHALL NOT quit the application and SHALL NOT alter the underlying view state. While the overlay is open, keys other than the dismiss keys SHALL NOT act on the underlying screen.

#### Scenario: Dismiss with `?`
- **WHEN** the overlay was opened from `ModeIndex` and the user presses `?`
- **THEN** the overlay closes and the mode is `ModeIndex` with its prior cursor and filter state intact

#### Scenario: Dismiss with `Esc`
- **WHEN** the overlay is open and the user presses `Esc`
- **THEN** the overlay closes and returns to the originating mode without quitting

#### Scenario: Dismiss with `q` does not quit
- **WHEN** the overlay is open and the user presses `q`
- **THEN** the overlay closes and the application keeps running in the originating mode

#### Scenario: Other keys are inert while open
- **WHEN** the overlay is open and the user presses a non-dismiss key such as `j`
- **THEN** the underlying screen does not scroll or move its cursor and the overlay remains open

### Requirement: Filter input takes precedence over `?`
When the index filter input is active (`ModeIndex` with the filter capturing text), pressing `?` SHALL be treated as filter text input rather than opening the overlay, so a `?` character can be entered into a filter query.

#### Scenario: `?` types into an active filter
- **WHEN** the mode is `ModeIndex`, the filter input is active, and the user presses `?`
- **THEN** the overlay does not open and `?` is appended to the filter text

### Requirement: Help-bar advertises the `?` affordance
The TUI SHALL include a `?: help` affordance in the help bar so the keyboard-help overlay is discoverable from the regular screens. The affordance SHALL NOT be shown in states where `?` is not available, such as while the index filter input is active.

#### Scenario: Help bar shows the help affordance
- **WHEN** the help bar is rendered for a mode where `?` opens the overlay
- **THEN** the help bar text includes a `?: help` affordance

#### Scenario: Help affordance hidden during filter input
- **WHEN** the mode is `ModeIndex` and the filter input is active
- **THEN** the help bar shows the filter prompt and does not advertise `?: help`
