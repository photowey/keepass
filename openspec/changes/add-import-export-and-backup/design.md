## Context

The project already has a versioned encrypted vault and a local-first operating model, but it still relies on ad hoc file copying or manual recreation for migration and recovery. Import/export and backup/restore are related but distinct:

- export/import moves logical entry data
- backup/restore preserves operational state, including encrypted vault and config

They should share file-handling utilities where possible, but not collapse into a single opaque workflow.

## Goals / Non-Goals

**Goals:**

- Provide a stable JSON export format for logical entry transfer.
- Support importing that JSON into an existing vault with explicit conflict strategy.
- Provide a backup bundle that captures vault and config files together.
- Provide restore behavior into a target keepass home directory with explicit overwrite rules.
- Keep all workflows explicit and operator-controlled.

**Non-Goals:**

- Remote sync
- Incremental backups
- Cross-format migration from third-party password managers in this change
- Automatic background snapshots

## Decisions

### 1. Separate “data transfer” from “backup/restore”

Decision:

- Export/import handles logical entries as plaintext JSON after unlock.
- Backup/restore handles encrypted vault/config files as operational assets.

Why:

- data portability and disaster recovery have different trust models
- mixing them makes operator intent ambiguous

Alternative considered:

- one combined archive format for everything
  - rejected because it makes selective migration harder and obscures trust boundaries

### 2. Use JSON as the v1 transfer format

Decision:

- Define a versioned JSON export schema with document metadata and entry list.

Why:

- easy to test
- easy to inspect
- stable for scripting

Alternative considered:

- CSV export/import
  - rejected because entry structure is richer than flat tabular data

### 3. Require explicit conflict strategy during import

Decision:

- Import must expose explicit conflict handling such as `fail`, `skip`, or `overwrite`.

Why:

- hidden overwrite behavior is dangerous for secrets
- import semantics must be deterministic

Alternative considered:

- default overwrite on alias conflict
  - rejected because it is too destructive

### 4. Use a directory-oriented backup layout

Decision:

- A backup operation creates a timestamped directory bundle containing:
  - encrypted vault
  - config file
  - manifest metadata

Why:

- easier to inspect and restore
- avoids archive-tool dependence in the core logic

Alternative considered:

- produce only zip/tar archives
  - rejected because core restore logic becomes more format-coupled

## Risks / Trade-offs

- [Plaintext export files are sensitive] → Mitigation: require explicit operator intent and document exposure clearly.
- [Import conflict rules can still surprise users] → Mitigation: make strategy explicit and default to fail-safe behavior.
- [Restore can overwrite healthy state] → Mitigation: require explicit overwrite behavior for destructive restore paths.
- [Backup format may need evolution] → Mitigation: include manifest versioning from the start.

## Migration Plan

1. Add transfer schema and import/export commands.
2. Add backup manifest and restore command.
3. Document plaintext export risk and backup restore workflow.
4. Rollback is straightforward because new commands are additive and do not change vault format.

## Open Questions

- Whether future versions should add archive packaging on top of the directory bundle.
- Whether future versions should support third-party import adapters as separate changes.
