## 1. Vault Metadata Inspection

- [x] 1.1 Add a vault inspection API that reads format and KDF metadata without requiring decryption
- [x] 1.2 Add tests for supported and invalid vault metadata inspection paths

## 2. Audit Report

- [x] 2.1 Introduce a structured audit report that evaluates config presence, vault presence, permissions, and rehash recommendation
- [x] 2.2 Add tests for healthy and mismatched config/vault audit scenarios

## 3. Doctor Command

- [x] 3.1 Add a `doctor` command with text and JSON output
- [x] 3.2 Update README and design/operations docs to describe the audit workflow
