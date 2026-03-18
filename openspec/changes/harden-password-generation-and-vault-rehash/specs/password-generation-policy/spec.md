## ADDED Requirements

### Requirement: Password generator SHALL support named presets
The system SHALL support named password generator presets so users can select a predefined alphabet policy without manually editing a raw alphabet string.

#### Scenario: Default config uses the compatibility preset
- **WHEN** the system creates a new default config
- **THEN** `password_generator.preset` SHALL be set to `compatible`

#### Scenario: Known preset resolves to a concrete alphabet
- **WHEN** password generation is requested and `password_generator.alphabet` is empty
- **THEN** the system SHALL resolve the effective alphabet from the configured preset

### Requirement: Explicit alphabet SHALL override preset selection
The system SHALL treat `password_generator.alphabet` as an explicit override when it is non-empty.

#### Scenario: Custom alphabet is configured
- **WHEN** `password_generator.alphabet` contains a non-empty value
- **THEN** password generation SHALL use that alphabet instead of any preset value

### Requirement: Unsupported preset values SHALL be rejected
The system SHALL reject invalid password generator preset values during config validation.

#### Scenario: Unknown preset is loaded from config
- **WHEN** config loading encounters an unsupported `password_generator.preset`
- **THEN** config validation SHALL fail with an explicit error
