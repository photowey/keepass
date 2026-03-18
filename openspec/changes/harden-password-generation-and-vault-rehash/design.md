## Context

The current code line already contains three related hardening moves:

- stronger default Argon2id memory settings
- password generator preset support with optional custom alphabet override
- a `rehash` command that rewrites the vault with the current configured KDF parameters

These changes affect multiple modules:

- `configs`
- `internal/password`
- `internal/manager`
- `cmd/cmder/root`
- `cmd/cmder/rehash`
- tests and user-facing docs

They belong to a single optimization track because they address the same operational problem: users need a practical way to strengthen password-related defaults over time without manually rewriting vault internals.

## Goals / Non-Goals

**Goals:**

- Provide a stronger default offline-resistance baseline for newly created configs.
- Replace raw-alphabet-only configuration with a more maintainable preset model.
- Preserve advanced custom alphabet support for environments with special constraints.
- Add a maintenance workflow for re-encrypting existing vaults after config hardening.
- Keep the change backward compatible for existing config files that still specify an explicit alphabet.

**Non-Goals:**

- Changing the vault file format version.
- Rotating or changing the master password during rehash.
- Introducing remote policy management or profile synchronization.
- Automatically rewriting user config files in place just because new defaults were introduced.

## Decisions

### 1. Keep password generation as “preset first, custom alphabet override second”

Decision:

- The config model adds `password_generator.preset`.
- `password_generator.alphabet` remains supported as an explicit override.
- Runtime behavior resolves the effective alphabet as:
  1. use `alphabet` if present
  2. otherwise use `preset`
  3. if `preset` is empty, fall back to `compatible`

Why:

- presets give users safe defaults without forcing them to handcraft character sets
- explicit alphabet override preserves full control for advanced use cases
- this allows backward compatibility because existing configs with `alphabet` continue to work

Alternative considered:

- Store only raw alphabets and document recommended values
  - rejected because it keeps policy definition too low-level and too easy to drift

### 2. Keep the default preset compatibility-oriented

Decision:

- The default preset remains `compatible`, not symbol-heavy by default.

Why:

- many password fields, shell contexts, and manual entry flows are still sensitive to broad symbol sets
- the project already gains substantial security from password length, alphabet size, and stronger Argon2 defaults
- users who need more symbols can switch to `symbols` or `strict-high-entropy`

Alternative considered:

- Default to a symbol-heavy alphabet
  - rejected because compatibility breakage is more immediate and common than entropy shortfall for the current length

### 3. Implement vault maintenance as a dedicated `rehash` command

Decision:

- Add `keepass rehash` as a first-class maintenance command.
- The command unlocks the vault with the current master password, loads the document, and saves it again using the current config values.

Why:

- it gives users an explicit workflow after hardening Argon2 parameters
- it reuses the normal encode/save path instead of inventing a second rewrite path
- it guarantees fresh salt and nonce generation through the existing codec implementation

Alternative considered:

- Auto-rewrite the vault during any successful unlock
  - rejected because it introduces hidden writes and surprising side effects

### 4. Do not combine rehash with master password rotation

Decision:

- `rehash` preserves the current master password.

Why:

- changing KDF cost and changing master password are separate operator concerns
- keeping rehash narrow reduces ambiguity and failure surface

Alternative considered:

- allow password change inside `rehash`
  - rejected because it conflates migration and credential rotation

## Risks / Trade-offs

- [Users may expect stronger defaults to automatically protect old vaults] → Mitigation: provide an explicit `rehash` workflow and document it clearly.
- [A preset model can hide the actual alphabet from advanced users] → Mitigation: keep `alphabet` override support and document preset intent.
- [Higher Argon2 memory increases unlock cost on weaker machines] → Mitigation: use rehash as an operator-controlled step rather than a forced migration.
- [Users may assume `symbols` is always better than `compatible`] → Mitigation: document that preset choice is a compatibility vs complexity trade-off, not a universal ranking.

## Migration Plan

1. New configs use the stronger Argon2 default and the `compatible` preset.
2. Existing configs with explicit `alphabet` continue to load without modification.
3. Users who increase Argon2 settings in config run `keepass rehash` to rewrite existing vaults.
4. Rollback remains straightforward:
   - restore the previous config values
   - run `keepass rehash` again if a rollback of vault KDF settings is required

## Open Questions

- Whether a future release should expose a command to print the effective alphabet resolved from `preset` plus optional `alphabet`.
- Whether future policy work should add site-specific password profiles instead of a single global preset.
