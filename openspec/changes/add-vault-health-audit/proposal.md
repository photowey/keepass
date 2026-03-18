## Why

The project now has stronger password-hardening defaults and a `rehash` workflow, but users still lack a single command that explains their current vault health. Operators need a non-destructive way to inspect whether configuration, file permissions, and stored KDF parameters are aligned before and after maintenance changes.

## What Changes

- Add a `doctor` command for vault health auditing.
- Report config presence, vault presence, resolved paths, effective password generation policy, and configured Argon2 settings.
- Inspect the vault file header without requiring the master password to surface stored format and KDF metadata.
- Detect whether the vault KDF parameters differ from the current config and recommend `rehash` when needed.
- Check restrictive filesystem permissions on supported platforms and surface warnings when local files are weaker than expected.
- Support both text and JSON output for operator and automation use.

## Capabilities

### New Capabilities
- `vault-health-audit`: Non-destructive inspection of local config, vault metadata, filesystem posture, and maintenance recommendations.

### Modified Capabilities

None.

## Impact

- Affected code:
  - `cmd/cmder/root`
  - new `cmd/cmder/doctor`
  - new `internal/audit`
  - `internal/vault`
  - `cmd/cmder/common`
- Affected user-facing behavior:
  - new maintenance command
  - new machine-readable audit output
- Affected docs:
  - README files
  - design and operations docs
