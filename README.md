# `keepass`

English | [中文](README.zh-CN.md)

`keepass` is a local-first CLI password manager written in Go.

It is designed around three goals:

- `Security first`: one master password unlocks a versioned encrypted vault file.
- `Low-friction`: short commands, interactive prompts, and safe defaults.
- `Fast lookup`: exact alias match first, then unique prefix match.

## What It Stores

Each entry is identified by a unique `alias` and can contain:

- `username`
- `password`
- `uri`
- `note`
- `tags`

Example:

- `github` -> `hellopass`
- `gitea` -> `hellopass`

The usernames can be the same. The `alias` is the unique lookup key.

## Security Model

- Vault file: `~/.keepass/keepass.kp`
- Config file: `~/.keepass/keepass.config.json`
- The config file stores only non-sensitive settings.
- The vault file stores encrypted entry data.
- The vault format includes a mandatory `format_version`.
- Parsing is strict: unknown versions fail closed.
- The master password is required to initialize and unlock the vault.
- Passwords are never stored in plaintext on disk.

## Quick Start

Initialize the vault:

```bash
keepass init
```

Add entries:

```bash
keepass add github hellopass --uri https://github.com --note "personal" --tag code
keepass add gitea hellopass --uri https://gitea.example.com --note "work" --tag code
```

When adding an entry:

- If you type an account password, the CLI asks for confirmation.
- If you leave it blank, `keepass` generates one for you.

List entries:

```bash
keepass list
keepass list --tag code
keepass list --json
```

Get a summary:

```bash
keepass get github
keepass get gith
```

Reveal the password explicitly:

```bash
keepass get gith --reveal
keepass get gith --json
keepass get gith --json --reveal
keepass get gith --copy
keepass get gith --copy --copy-timeout 0
```

Update and delete:

```bash
keepass update github
keepass delete github
```

Transfer and recovery:

```bash
keepass export --path ./entries.json
keepass import --path ./entries.json --conflict overwrite
keepass backup --path ./backup-bundle
keepass restore --path ./backup-bundle --force
```

Credential hygiene:

```bash
keepass audit --json
keepass rotate github --generate
```

Rewrite the vault using the current Argon2 settings:

```bash
keepass rehash
```

Inspect effective config:

```bash
keepass config
keepass config --json
```

Audit local vault health:

```bash
keepass doctor
keepass doctor --json
```

## Automation Notes (Non-interactive Mode)

When stdin is not a TTY (scripts, CI, pipes), some commands avoid prompting to prevent hanging:

- `keepass add` requires `alias` and `username` as arguments (no interactive prompts).
- `keepass update` requires explicit mutation flags such as `--username`, `--password`, `--clear-uri`, or `--clear-note`.
- `keepass delete` requires `--yes` to skip confirmation in non-interactive mode.
- You can force this behavior even in a TTY with `--non-interactive`.

## Exit Codes

- `1`: generic error
- `2`: usage / invalid arguments
- `3`: not initialized (missing config/vault)
- `4`: unlock failed (wrong master password)

## Shell Completion

Generate completion scripts:

```bash
keepass completion bash
keepass completion zsh
keepass completion fish
keepass completion powershell
```

## Alias Resolution

Lookup rules are:

1. Exact alias match
2. Unique prefix match
3. Ambiguous prefix -> fail with all candidates

Examples:

- `keepass get github` -> exact match
- `keepass get gith` -> unique prefix match
- `keepass get gi` -> fails if both `github` and `gitea` exist

## Password Generation

If you do not provide an account password, `keepass` generates one with secure randomness.

The generator accepts any non-empty alphabet from config.

The default alphabet intentionally avoids most special symbols for compatibility across websites, shells, and manual entry.
If your environment requires additional symbols, switch `password_generator.preset` or set a custom `password_generator.alphabet`.

Built-in presets:

- `compatible`
  - default, optimized for broad website and shell compatibility
- `symbols`
  - adds a moderate set of special symbols
- `strict-high-entropy`
  - uses a larger mixed alphabet with more symbols

Default settings live in `~/.keepass/keepass.config.json`:

```json
{
  "version": 1,
  "vault": {
    "path": "~/.keepass/keepass.kp",
    "format_version": 1
  },
  "security": {
    "argon2id": {
      "time": 3,
      "memory_kib": 262144,
      "threads": 4,
      "key_length": 32
    }
  },
  "password_generator": {
    "default_length": 21,
    "preset": "compatible"
  }
}
```

## Testing

The project includes:

- config validation tests
- password generator tests
- vault format and crypto tests
- manager rule tests
- command flow integration tests
- a fuzz entry point for vault decoding

Run:

```bash
GOCACHE=/tmp/go-cache go test ./...
```

## Release Integrity

GitHub Release artifacts include:

- per-file SHA256 checksums in `SHA256SUMS.txt`
- GitHub artifact attestations for build provenance
