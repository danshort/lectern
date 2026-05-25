## ADDED Requirements

### Requirement: Ordenar cambios activos por fecha de creación
`LoadFrom()` SHALL sort active changes by their `created` date in descending order (newest first). Changes without a `created` date SHALL be placed after all dated changes and sorted alphabetically by name (stable sort as tiebreaker). Changes with equal `created` dates SHALL keep their relative input order (stable sort).

#### Scenario: Cambios con fechas distintas
- **WHEN** two active changes have `created` dates `2026-05-01` and `2026-05-10`
- **THEN** the change from `2026-05-10` appears before the change from `2026-05-01` in the returned list

#### Scenario: Cambios sin fecha van al final
- **WHEN** one active change has no `created` date and another has `2026-05-01`
- **THEN** the dated change appears first and the undated change appears after, sorted alphabetically

#### Scenario: Cambios con la misma fecha
- **WHEN** two changes have the same `created` date
- **THEN** their relative order is stable (preserved from directory listing)
