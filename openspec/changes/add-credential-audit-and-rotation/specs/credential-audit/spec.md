## ADDED Requirements

### Requirement: User shall be able to audit unlocked credential hygiene
The system SHALL support auditing unlocked entries for credential hygiene findings.

#### Scenario: Audit reports stale passwords
- **WHEN** an entry password age exceeds the requested audit threshold
- **THEN** the audit output SHALL report the entry as stale

#### Scenario: Audit reports duplicate passwords
- **WHEN** two or more entries share the same password
- **THEN** the audit output SHALL report those entries as a duplicate-password finding

### Requirement: Audit shall report missing metadata findings
The system SHALL report entries with missing operational metadata such as username or URI.

#### Scenario: Entry is missing URI
- **WHEN** an entry has no URI
- **THEN** the audit output SHALL include a missing-metadata finding for that alias
