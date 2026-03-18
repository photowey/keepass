## ADDED Requirements

### Requirement: User shall be able to rotate a target entry password
The system SHALL support rotating the password for a target alias through a dedicated command.

#### Scenario: Rotate with generated password
- **WHEN** the user runs the rotate command with generate mode
- **THEN** the system SHALL generate a new password and store it on the target entry

#### Scenario: Rotate with manual password
- **WHEN** the user runs the rotate command with an explicit password value
- **THEN** the system SHALL store the provided password on the target entry

### Requirement: Rotation output shall not expose the new password by default
The rotation workflow SHALL avoid printing the new password unless the user explicitly requests one-time disclosure.

#### Scenario: Rotation completes without reveal
- **WHEN** the user runs a successful rotate command without reveal or copy flags
- **THEN** the system SHALL confirm rotation without printing the plaintext password
