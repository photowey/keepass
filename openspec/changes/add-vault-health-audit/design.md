## Context

The codebase now includes stronger password hardening defaults and a `rehash` workflow, but there is still no single operator-facing command that explains the current local security posture. Users can inspect config manually, but they cannot easily answer questions such as:

- does the vault exist?
- what KDF parameters were actually used for the current vault file?
- do local file permissions match the intended model?
- does the current config differ from the vault’s stored KDF metadata?
- is `rehash` recommended?

The new change needs to be non-destructive, automation-friendly, and consistent with the project’s local-first operating model.

## Goals / Non-Goals

**Goals:**

- Add a `doctor` command that audits local keepass state without requiring the master password.
- Surface both configured Argon2 settings and vault-stored Argon2 settings when available.
- Detect mismatches between current config and vault header metadata and recommend `rehash`.
- Report filesystem permission health where the platform provides meaningful local mode bits.
- Support both text output and JSON output for operators and scripts.

**Non-Goals:**

- Attempt to unlock the vault or validate the master password.
- Mutate config or vault files automatically.
- Repair incorrect permissions automatically.
- Replace `config` or `rehash`; this command is diagnostic, not mutating.

## Decisions

### 1. Add a dedicated `internal/audit` layer

Decision:

- Introduce a small audit/reporting package rather than embedding health logic directly inside the command.

Why:

- keeps `cmd/cmder/doctor` thin
- allows tests to validate health logic without shell wiring
- centralizes status aggregation in one report model

Alternative considered:

- implement all checks directly inside the command package
  - rejected because it would make diagnostics harder to test and easier to drift

### 2. Inspect vault metadata without decrypting the vault

Decision:

- Add an internal vault inspection function that reads file magic, format version, header length, and v1 KDF metadata without requiring the master password.

Why:

- `doctor` should be safe to run without secrets
- the vault header already contains the data needed for KDF posture checks
- this avoids mixing diagnostics with unlock logic

Alternative considered:

- require a master password and inspect by decrypting
  - rejected because it adds friction and weakens the usefulness of a passive health command

### 3. Use a structured report with machine-readable severity

Decision:

- The audit result should contain structured checks with statuses such as `ok`, `warn`, and `error`.

Why:

- text output remains readable for humans
- JSON output becomes stable for automation
- severity helps downstream tooling make decisions without scraping prose

Alternative considered:

- return only a free-form summary string
  - rejected because it is weak for automation and hard to extend

### 4. Treat permission checks as platform-aware, not universal

Decision:

- Permission checks are authoritative on Unix-like systems and best-effort or skipped where the model is not meaningful.

Why:

- file mode expectations differ across platforms
- a false alarm is worse than a skipped check for this type of command

Alternative considered:

- force the same permission rule on every platform
  - rejected because it would produce misleading health output

## Risks / Trade-offs

- [Vault inspection logic duplicates some header parsing concerns] → Mitigation: keep inspection focused on metadata and reuse existing format semantics.
- [Diagnostic output may grow over time and drift] → Mitigation: use a structured report model with explicit fields and statuses.
- [Permission warnings may still be environment-sensitive] → Mitigation: scope them clearly and skip unsupported cases rather than guessing.
- [Users may confuse `doctor` with repair behavior] → Mitigation: keep recommendations explicit and leave mutation to `rehash` or manual fixes.

## Migration Plan

1. Add the new non-destructive command and report model.
2. Document how `doctor` relates to `config` and `rehash`.
3. Rollback is simple:
   - remove the command registration
   - remove the audit/report package
   - remove vault metadata inspection helpers if needed

## Open Questions

- Whether a future version should support a machine-stable exit code policy for audit severity.
- Whether `doctor` should later expose a `--quiet` or `--fail-on-warn` automation mode.
