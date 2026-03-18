## ADDED Requirements

### Requirement: Update shall be automation-safe in non-interactive mode
The system SHALL ensure that `update` does not prompt for field values in non-interactive mode.

#### Scenario: Non-interactive update with explicit flags
- **WHEN** the user runs `keepass update <alias>` in non-interactive mode with supported mutation flags
- **THEN** the system SHALL apply the requested changes without prompting for additional field input

#### Scenario: Non-interactive update without mutation intent
- **WHEN** the user runs `keepass update <alias>` in non-interactive mode without any mutation flags
- **THEN** the system SHALL fail with a usage-oriented error

### Requirement: Update shall support explicit scripted clearing of optional fields
The system SHALL support explicit clearing of optional fields in non-interactive mode.

#### Scenario: Clear URI and note through flags
- **WHEN** the user runs `keepass update <alias> --clear-uri --clear-note`
- **THEN** the system SHALL clear the stored URI and note values without prompting

### Requirement: Update shall support manual password replacement in non-interactive mode
The system SHALL allow password replacement through an explicit flag in non-interactive mode.

#### Scenario: Replace password through flag
- **WHEN** the user runs `keepass update <alias> --password <value>` in non-interactive mode
- **THEN** the system SHALL update the password to the provided value without interactive prompts

### Requirement: Delete shall require explicit confirmation in non-interactive mode
The system SHALL fail fast when a non-interactive delete is attempted without explicit confirmation intent.

#### Scenario: Non-interactive delete without yes flag
- **WHEN** the user runs `keepass delete <alias>` in non-interactive mode without `--yes`
- **THEN** the system SHALL fail with a usage-oriented error

#### Scenario: Non-interactive delete with yes flag
- **WHEN** the user runs `keepass delete <alias> --yes` in non-interactive mode
- **THEN** the system SHALL delete the entry without prompting for delete confirmation
