## Context

The codebase already introduced a global `--non-interactive` flag and uses it effectively in `add`, but mutation commands are still uneven:

- `update` is interactive-first and can still block or behave ambiguously in automation
- `delete` only becomes deterministic in scripts when the caller already knows to pass `--yes`
- several automation-contract violations currently surface as generic errors instead of usage-oriented ones

This makes the CLI less predictable than it should be for CI, shell scripts, or other tooling.

## Goals / Non-Goals

**Goals:**

- Make `update` safe and explicit in non-interactive mode.
- Require explicit confirmation flags for destructive scripted operations.
- Return usage-style failures for missing required mutation intent.
- Preserve the existing interactive UX for local terminal usage.

**Non-Goals:**

- Converting every command to a JSON-first mutation interface.
- Introducing config-file-based batch mutation workflows.
- Changing master password prompt behavior.

## Decisions

### 1. Make non-interactive `update` flag-driven

Decision:

- In non-interactive mode, `update` must not prompt.
- All mutations must come from explicit flags.
- If no mutation flags are provided, the command fails with a usage-style error.

Why:

- scripts need deterministic behavior
- hidden prompts are unacceptable in automation contracts
- update intent should be explicit

Alternative considered:

- keep prompt fallbacks even in non-interactive mode
  - rejected because it defeats the purpose of a non-interactive contract

### 2. Add explicit clear flags for optional fields

Decision:

- Add `--clear-uri` and `--clear-note`.

Why:

- scripting must support both replacement and clearing
- relying on interactive blank-input semantics is not automation-safe

Alternative considered:

- overload empty strings on `--uri` and `--note`
  - rejected because empty-string flag values are ambiguous and shell-hostile

### 3. Add explicit `--password` for scripted update

Decision:

- `update` gets a `--password` flag for non-interactive manual replacement.

Why:

- generated-password mode already exists through `--generate`
- scripts still need deterministic manual password replacement

Alternative considered:

- require only `--generate` in non-interactive mode
  - rejected because it removes a valid automation use case

### 4. Require `--yes` for non-interactive delete

Decision:

- `delete` in non-interactive mode fails fast unless `--yes` is provided.

Why:

- silent cancellation is a poor automation contract
- destructive commands should require explicit operator intent in scripts

Alternative considered:

- keep current cancellation behavior
  - rejected because scripts can misread it as success without mutation intent

## Risks / Trade-offs

- [More flags increase CLI surface area] → Mitigation: keep flags scoped to clear automation needs only.
- [Users may expect blank flag values to clear fields] → Mitigation: document explicit clear flags instead of implicit blank-string conventions.
- [Usage errors require careful mapping] → Mitigation: centralize a small helper for usage-oriented command failures.

## Migration Plan

1. Add flag-driven non-interactive mutation behavior.
2. Add tests for non-interactive update/delete and built-binary exit semantics.
3. Update docs to describe the stricter automation contract.
4. Rollback is straightforward because changes are isolated to command behavior and CLI helpers.

## Open Questions

- Whether future versions should add a dedicated `--stdin-password` pattern for automation environments that avoid shell history.
