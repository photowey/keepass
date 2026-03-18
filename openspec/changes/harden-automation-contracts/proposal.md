## Why

The current CLI is strong for interactive local use, but mutation commands are still inconsistent for automation. `add` has explicit non-interactive behavior, while `update` still prompts by default and `delete` can silently cancel in non-interactive contexts instead of failing fast with a usage-oriented contract.

## What Changes

- Harden mutation command behavior for non-interactive environments.
- Make `update` fully usable from scripts through explicit flags instead of prompt-only mutation paths.
- Add explicit clearing flags for optional fields that currently require interactive overwrite behavior.
- Require `delete --yes` in non-interactive mode so scripts fail fast instead of silently cancelling.
- Return usage-style exit semantics for automation-contract violations.

## Capabilities

### New Capabilities
- `automation-safe-mutations`: Stable non-interactive behavior for update and delete operations, including explicit mutation flags and fail-fast usage errors.

### Modified Capabilities

None.

## Impact

- Affected code:
  - `cmd/cmder/update`
  - `cmd/cmder/delete`
  - `cmd/cmder/common`
  - `cmd/cmder/root`
  - tests for command execution and exit codes
- Affected user-facing behavior:
  - non-interactive update and delete semantics
  - new flags for scripted updates
- Affected docs:
  - README files
  - CLI and design docs
