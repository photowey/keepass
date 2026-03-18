## Context

The current CLI exposes entry CRUD and a vault health audit, but it does not provide an operator-focused view of credential quality inside the decrypted dataset. Password age, password reuse, and missing metadata are all useful hygiene signals, and rotation deserves its own focused workflow instead of forcing users to repurpose the generic `update` command.

## Goals / Non-Goals

**Goals:**

- Provide a credential audit command over unlocked entries.
- Surface stale-password, duplicate-password, and missing-metadata findings.
- Add a dedicated rotate command that updates the target entry password explicitly.
- Keep audit logic and rotation logic scriptable and deterministic.

**Non-Goals:**

- Guess password entropy from external breach corpora
- Automatic scheduled rotation
- Account integration with remote services
- MFA/TOTP analysis

## Decisions

### 1. Make audit an unlocked logical analysis, not a file-level health check

Decision:

- Credential audit runs after vault unlock and inspects entry contents.

Why:

- the findings depend on logical entry values, not file metadata
- it complements rather than duplicates `doctor`

Alternative considered:

- merge credential audit into `doctor`
  - rejected because the trust model and required unlock behavior are different

### 2. Define rotation as a focused command instead of overloading `update`

Decision:

- Add a dedicated `rotate` command for password rotation.

Why:

- rotation is a common operator workflow with specific intent
- a dedicated command can support generate/manual modes and one-time disclosure options cleanly

Alternative considered:

- keep using `update --generate`
  - rejected because it hides the operator intent and makes future rotation-specific UX harder

### 3. Keep duplicate-password detection local and exact

Decision:

- Detect duplicate passwords by exact equality across unlocked entries.

Why:

- simple and deterministic
- no external dependency

Alternative considered:

- attempt fuzzy or breach-style password quality analysis
  - rejected because it expands scope too aggressively

## Risks / Trade-offs

- [Credential audit requires plaintext access after unlock] → Mitigation: keep output controlled and avoid printing passwords in audit reports.
- [Rotation command may overlap with update semantics] → Mitigation: define it as the dedicated password-only workflow and keep update as general metadata mutation.
- [Duplicate detection can be expensive for large vaults] → Mitigation: accept simple in-memory analysis for the current project scale.

## Migration Plan

1. Add credential audit report model and command.
2. Add rotation command and manager support.
3. Document audit findings and rotation workflow.
4. Rollback is additive and isolated to new commands and analysis logic.

## Open Questions

- Whether future versions should support audit severity thresholds as exit codes.
- Whether rotation should later support batch workflows by tag or query.
