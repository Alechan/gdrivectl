---
name: gdrivectl-drive-ops
description: Use this skill when a task needs Google Drive/Docs operations through gdrivectl (doctor, search, file metadata, doc tabs, and doc export), including auth/scope/network remediation and exit-code-aware handling.
---

# gdrivectl Drive Ops

Use this skill for Drive/Docs reads via `gdrivectl`.

## When to use

- User asks to inspect Drive files, search documents, check file metadata, inspect Google Docs tabs, or export Docs content.
- User needs structured troubleshooting for `gcloud` auth/scope/network/API failures from this CLI.

## When not to use

- Tasks requiring write/mutation operations in Drive/Docs (not supported by this CLI v0.1).
- Tasks unrelated to Google Drive/Docs.

## Preconditions

- `gcloud` installed and authenticated.
- prefer installed binary mode:
  - `command -v gdrivectl`
- If scope issues appear, run:
  - `gcloud auth login --enable-gdrive-access --update-adc`

## Command map

- Invocation mode:
  - binary-first: `gdrivectl ...`
  - source fallback (contributors): `go run ./cmd/gdrivectl ...` from gdrivectl repo root
- Environment preflight: `gdrivectl doctor [--json]`
- Drive search: `gdrivectl search --query <q> [--page-size N] [--json]`
- File metadata: `gdrivectl file-meta --id <FILE_ID> [--json]`
- Docs tabs: `gdrivectl doc-tabs --id <DOC_ID> [--json]`
- Doc export: `gdrivectl doc-export --id <DOC_ID> --mime <MIME> [--out <PATH>]`

Use `--json` by default for structured commands. For `doctor`, prefer `--json` when results are consumed programmatically.

## Standard flow

1. Confirm task intent and required IDs/filters.
2. Run command directly if preconditions are known-good.
3. If command fails and environment is unknown, run `doctor` before retrying.
4. Apply remediation based on exit code and error category.

## Exit-code remediation

- `2` (`validation/config`): fix flags or `--gcloud-bin` path.
- `3` (`auth/scope`): re-authenticate and ensure Drive scope.
  - Run: `gcloud auth login --enable-gdrive-access --update-adc`
  - If it persists, run: `gcloud auth print-access-token`
  - Optional diagnostic: `gcloud auth application-default print-access-token`
  - If token diagnostics still fail due to sandbox/network/config-store constraints, rerun the target `gdrivectl` command with escalated/unsandboxed execution.
- `4` (`network`): retry with larger timeout (for example `--timeout 60s`) and verify connectivity.
- `5` (`api`): validate IDs/query/mime/access.

## Safety rules

- Default to read-only operations.
- Never invent file IDs or document IDs.
- If exporting to disk, use explicit user-provided or clearly scoped output paths.
- Do not claim success without command evidence.

## Response style

- Provide exact command(s) run.
- State invocation mode used (`binary` or `source fallback`).
- Summarize key output fields only.
- Include next remediation step when command fails.
- For escalated retries, always report:
  - original command
  - failing exit code/category
  - escalated retry command
  - final outcome and output path (for `doc-export`)
