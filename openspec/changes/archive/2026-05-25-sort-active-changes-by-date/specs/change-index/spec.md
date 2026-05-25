## ADDED Requirements

### Requirement: Orden de cambios activos por fecha
Active changes in the index SHALL be displayed in creation date order, newest first, as provided by the loader.

#### Scenario: Índice con cambios de fechas variadas
- **WHEN** the index is rendered and active changes have different creation dates
- **THEN** the newest change appears first in the "Active Changes" section

#### Scenario: Cambio sin fecha aparece al final
- **WHEN** an active change has no `created` date
- **THEN** it appears after all dated changes in the "Active Changes" section
