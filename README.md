# `keepass`

English | [中文](./README_zh_CN.md)

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

- `github` -> `abc`
- `gitea` -> `abc`

The usernames can be the same. The `alias` is the unique lookup key.

## Security Model

- Vault file: `~/.keepass/default.kp`
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
keepass add github abc --uri https://github.com --note "personal" --tag code
keepass add gitea abc --uri https://gitea.example.com --note "work" --tag code
```

When adding an entry:

- If you type an account password, the CLI asks for confirmation.
- If you leave it blank, `keepass` generates one for you.

List entries:

```bash
keepass list
keepass list --tag code
```

Get a summary:

```bash
keepass get github
keepass get gith
```

Reveal the password explicitly:

```bash
keepass get gith --reveal
```

Update and delete:

```bash
keepass update github
keepass delete github
```

Inspect effective config:

```bash
keepass config
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

Default settings live in `~/.keepass/keepass.config.json`:

```json
{
  "version": 1,
  "vault": {
    "path": "~/.keepass/default.kp",
    "format_version": 1
  },
  "security": {
    "argon2id": {
      "time": 3,
      "memory_kib": 65536,
      "threads": 4,
      "key_length": 32
    }
  },
  "password_generator": {
    "default_length": 21,
    "alphabet": "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz23456789-_"
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
