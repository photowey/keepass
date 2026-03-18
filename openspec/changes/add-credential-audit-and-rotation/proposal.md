## Why

The project now has stronger KDF posture, password generation presets, and a vault health audit, but it still does not help users inspect credential quality inside the vault itself. Operators need a planned workflow for identifying stale or risky entries and rotating passwords intentionally rather than manually discovering and fixing issues one by one.

## What Changes

- Add credential audit capabilities for local entry hygiene.
- Detect stale passwords by age threshold using `password_updated_at`.
- Detect duplicate passwords across aliases after unlock.
- Detect entries with missing or weak metadata such as missing username or URI.
- Add a dedicated rotation workflow that can generate or set a new password for a target alias and optionally print/copy the new password once.

## Capabilities

### New Capabilities
- `credential-audit`: Analyze unlocked entries for stale, duplicated, or incomplete credential posture.
- `credential-rotation`: Rotate a target entry password through a dedicated command and explicit output controls.

### Modified Capabilities

None.

## Impact

- Affected code:
  - `cmd/cmder/*`
  - `internal/manager`
  - new credential audit package
- Affected user-facing behavior:
  - new audit and rotate commands
  - new audit output and filtering semantics
- Affected docs:
  - README files
  - design and operations docs
