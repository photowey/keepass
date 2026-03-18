## Why

The current CLI manages secrets locally, but it still lacks first-class data portability and disaster-recovery workflows. Users need a planned, explicit way to export entries, import entries, and create encrypted local backups without manually copying internal files or reverse-engineering the vault layout.

## What Changes

- Add plaintext export of unlocked entries to a portable JSON format.
- Add import of exported JSON data back into the local vault with explicit conflict handling.
- Add encrypted backup creation for the current local state, including vault and config files.
- Add explicit restore support from a backup bundle into a target keepass home.
- Define a stable export schema and a stable backup layout for future compatibility.

## Capabilities

### New Capabilities
- `vault-data-transfer`: Export and import unlocked entry data through a stable JSON format with explicit conflict rules.
- `vault-backup-and-restore`: Create and restore local encrypted backups containing vault and config state.

### Modified Capabilities

None.

## Impact

- Affected code:
  - `cmd/cmder/*`
  - `internal/manager`
  - `internal/vault`
  - new import/export/backup support packages
- Affected user-facing behavior:
  - new transfer and recovery commands
  - new file formats and conflict-handling flags
- Affected docs:
  - README files
  - design and operations docs
