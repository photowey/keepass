## ADDED Requirements

### Requirement: Vault rehash SHALL rewrite the vault using current security parameters
The system SHALL provide a command that rewrites an existing vault using the current configured Argon2id parameters while preserving vault contents.

#### Scenario: Rehash is executed with a valid master password
- **WHEN** the user runs `keepass rehash` and provides the correct master password
- **THEN** the system SHALL load the current vault
- **THEN** the system SHALL save the same logical document using the current configured Argon2id parameters

### Requirement: Vault rehash SHALL preserve the current master password
The rehash workflow SHALL not change the user’s master password.

#### Scenario: Rehash completes successfully
- **WHEN** the user runs `keepass rehash`
- **THEN** subsequent unlock operations with the same master password SHALL continue to succeed

### Requirement: Vault rehash SHALL preserve existing entries
The rehash workflow SHALL preserve all stored entries and their readable content.

#### Scenario: Rehash is performed on a vault with entries
- **WHEN** the user runs `keepass rehash` on a non-empty vault
- **THEN** existing aliases, usernames, passwords, tags, and notes SHALL remain available after the rewrite
