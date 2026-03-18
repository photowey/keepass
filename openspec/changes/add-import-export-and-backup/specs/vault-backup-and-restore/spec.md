## ADDED Requirements

### Requirement: User shall be able to create a local backup bundle
The system SHALL support creating a backup bundle that contains the encrypted vault, config, and manifest metadata.

#### Scenario: Backup succeeds
- **WHEN** the user runs the backup command against an initialized keepass home
- **THEN** the system SHALL create a backup bundle containing the current vault and config files

### Requirement: User shall be able to restore from a backup bundle
The system SHALL support restoring vault and config state from a backup bundle into a target keepass home.

#### Scenario: Restore with explicit overwrite
- **WHEN** the user restores a backup into a target location with overwrite enabled
- **THEN** the system SHALL replace the target vault and config with the backup contents
