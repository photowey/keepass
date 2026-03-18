## 1. Security Baseline

- [x] 1.1 Increase the default Argon2id memory cost for newly created configurations
- [x] 1.2 Add validation and tests for the strengthened default security configuration

## 2. Password Generation Policy

- [x] 2.1 Introduce named password generator presets with a compatibility-focused default
- [x] 2.2 Preserve explicit custom alphabet override behavior and cover it with tests
- [x] 2.3 Update user-facing documentation to describe presets and when to use symbol-rich policies

## 3. Vault Rehash Workflow

- [x] 3.1 Add a `rehash` command that rewrites the vault using the current configured Argon2 settings
- [x] 3.2 Verify rehash preserves vault contents while updating the stored KDF parameters
