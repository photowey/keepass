## ADDED Requirements

### Requirement: Operator shall be able to audit local vault health without unlocking the vault
The system SHALL provide a non-destructive audit command that inspects local keepass state without requiring the master password.

#### Scenario: Doctor command runs on local state
- **WHEN** the user runs `keepass doctor`
- **THEN** the system SHALL inspect local config, vault metadata, and selected filesystem properties without requesting the master password

### Requirement: Audit output shall report config and vault security posture
The audit command SHALL report both configured security posture and vault-stored security metadata when available.

#### Scenario: Config and vault both exist
- **WHEN** the audit command runs with both config and vault files present
- **THEN** the output SHALL include configured Argon2 parameters
- **THEN** the output SHALL include vault format and vault-stored KDF metadata

### Requirement: Audit output shall detect when rehash is recommended
The audit command SHALL detect when the vault’s stored KDF parameters differ from the currently configured Argon2 settings.

#### Scenario: Config differs from stored vault KDF metadata
- **WHEN** the vault header Argon2 parameters are not equal to the current config Argon2 parameters
- **THEN** the audit output SHALL mark that rehash is recommended

### Requirement: Audit output shall support machine-readable format
The audit command SHALL support JSON output for automation.

#### Scenario: JSON mode is requested
- **WHEN** the user runs `keepass doctor --json`
- **THEN** the system SHALL emit a structured report containing audit statuses and recommendations
