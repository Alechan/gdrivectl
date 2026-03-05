# SDD-0007: Debug Playbook and Portability for gdrivectl

Status: Draft  
Version: 0.1  
Last updated: 2026-03-05

## 1. Purpose

Standardize a practical debugging playbook for `gdrivectl` that is:

- aligned with current CLI behavior
- portable across machines/shells
- clear about exit-code-based remediation

This SDD consolidates troubleshooting scenarios and review recommendations gathered during v0.1 hardening.

## 2. Baseline scenarios from original cheatsheet

1. `gdrivectl: command not found`
2. `gcloud binary not found` (exit `2`)
3. `unable to get access token` (exit `3`)
4. sandbox/permission errors with `~/.config/gcloud`
5. healthy `doctor --json` baseline
6. one-shot validation command

## 3. Additional recommendations (new)

### 3.1 Portability

- Avoid hardcoded absolute home paths in examples.
- Prefer `$HOME/...` or `$(command -v gcloud)` patterns.
- Keep absolute-path examples only as optional host-specific fallbacks.

### 3.2 Accuracy of expected output

- Healthy `doctor --json` examples should include current stable fields:
  - `gcloud_bin`
  - `gcloud_exists`
  - `token_ok`
  - `drive_ok`
  - `docs_ok`
  - optional `note`

### 3.3 Better diagnostics flow

- Before retrying any failed command:
  1. run `command -v gcloud`
  2. run `gcloud auth list`
  3. run `gdrivectl doctor --json`
- For network-related failures, explicitly test with larger timeout (`--timeout 60s`).
- For scope failures, always include exact remediation command:
  - `gcloud auth login --enable-gdrive-access --update-adc`

### 3.4 Sandbox-specific guidance

- Document that some environments block writes to default gcloud config paths.
- Add fallback pattern using explicit config dir when needed:
  - `CLOUDSDK_CONFIG=<writable_dir> gcloud ...`

### 3.5 Exit-code matrix visibility

- Every debug note should map symptoms to exit codes:
  - `2`: validation/config
  - `3`: auth/scope
  - `4`: network
  - `5`: API

### 3.6 Command style consistency

- Prefer `go run ./cmd/gdrivectl ...` for repo-local reproducibility in docs.
- Keep examples in JSON mode for structured commands and debugging automation.

## 4. Proposed canonical debug sequence

1. `go run ./cmd/gdrivectl --help`
2. `command -v gcloud`
3. `gcloud auth list`
4. `go run ./cmd/gdrivectl doctor --json --gcloud-bin \"$(command -v gcloud || echo gcloud)\"`
5. If failure code `3`: run scope re-auth command, retry step 4.
6. If failure code `4`: retry step 4 with `--timeout 60s`.
7. Then run target command (`search`, `file-meta`, `doc-tabs`, or `doc-export`).

## 5. Required documentation updates

- Keep debug playbook synchronized across:
  - `README.md` troubleshooting
  - `docs/TEST_PLAN.md` negative-path checks
  - any external disposable/internal runbooks

## 6. Acceptance criteria

- A new contributor can diagnose common failures with only documented steps.
- No critical debug examples rely on machine-specific absolute paths.
- Debug examples match current command output fields and exit-code contract.

## 7. Traceability

- Extends SDD-0004 documentation-completion requirements.
- Aligns with v0.1 output and error contracts in SDD-0003.
