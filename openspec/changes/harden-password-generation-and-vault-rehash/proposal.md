## Why

The current security-related iteration has improved several areas, but it was driven incrementally instead of as a single planned optimization track. Password hardening, password generator behavior, and vault maintenance need to be captured as one coherent change so the project has a stable baseline for future iterations.

## What Changes

- Increase the default Argon2id memory cost to strengthen offline resistance for newly created configurations.
- Introduce built-in password generator presets so users can choose between compatibility-oriented and symbol-rich generation without hand-editing raw alphabets.
- Keep support for explicit custom alphabets as an override for advanced environments.
- Add a `rehash` command that rewrites an existing vault using the current configured Argon2id parameters while preserving the current master password and stored entries.
- Align tests, README content, and design documents with the strengthened defaults and new vault maintenance workflow.

## Capabilities

### New Capabilities
- `password-generation-policy`: Configurable password generation presets with a compatibility-focused default and optional stronger symbol-rich presets.
- `vault-rehash`: A maintenance command that rewrites an existing vault using the current configured password-hardening parameters.

### Modified Capabilities

None.

## Impact

- Affected code:
  - `configs`
  - `internal/password`
  - `internal/manager`
  - `cmd/cmder/root`
  - `cmd/cmder/rehash`
- Affected user-facing behavior:
  - default config values
  - password generator configuration format
  - available CLI commands
- Affected documentation:
  - README files
  - design and operations docs
