## ADDED Requirements

### Requirement: User shall be able to export unlocked entry data to a stable JSON format
The system SHALL support exporting unlocked vault entries to a versioned JSON document.

#### Scenario: Export succeeds after unlock
- **WHEN** the user runs the export command with the correct master password
- **THEN** the system SHALL write a versioned JSON document containing the current logical entry data

### Requirement: User shall be able to import entry data with explicit conflict handling
The system SHALL support importing exported JSON data into an existing vault with explicit alias-conflict behavior.

#### Scenario: Import with fail-on-conflict
- **WHEN** the import command encounters an alias conflict and the strategy is `fail`
- **THEN** the system SHALL stop and report the conflict without partial overwrite

#### Scenario: Import with overwrite-on-conflict
- **WHEN** the import command encounters an alias conflict and the strategy is `overwrite`
- **THEN** the system SHALL replace the existing entry with the imported entry
