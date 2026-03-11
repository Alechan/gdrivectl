---
name: gdrivectl-drive-ops
description: Use this skill when a task needs Google Drive/Docs operations through gdrivectl (doctor, search, file metadata, doc tabs, doc export, and upload), including auth/scope/network remediation and exit-code-aware handling.
---

# gdrivectl Drive Ops

Use this skill for Drive/Docs operations via `gdrivectl`.

## When to use

- User asks to inspect Drive files, search documents, check file metadata, inspect Google Docs tabs, or export Docs content.
- User asks to upload a local file into Google Drive.
- User needs structured troubleshooting for `gcloud` auth/scope/network/API failures from this CLI.

## When not to use

- Destructive Drive operations (delete, permission mutation), which are not supported by this CLI.
- Tasks unrelated to Google Drive/Docs.

## Preconditions

- `gcloud` installed and authenticated.
- Prefer installed binary mode: `command -v gdrivectl`
- If scope issues appear, run: `gcloud auth login --enable-gdrive-access --update-adc`

## Command map

- Invocation mode:
  - binary-first: `gdrivectl ...`
  - source fallback (contributors): `go run ./cmd/gdrivectl ...` from gdrivectl repo root
- Environment preflight: `gdrivectl doctor [--json]`
- Drive search: `gdrivectl search --query <q> [--page-size N] [--json]`
- File metadata: `gdrivectl file-meta --id <FILE_ID> [--json]`
- Docs tabs: `gdrivectl doc-tabs --id <DOC_ID> [--json]`
- Doc export: `gdrivectl doc-export --id <DOC_ID> --mime <MIME> [--out <PATH>]`
- File upload: `gdrivectl upload --path <LOCAL_FILE> [--name <NAME>] [--parent-id <FOLDER_ID>] [--mime <MIME>] [--json]`

Use `--json` by default for structured commands (`search`, `file-meta`, `doc-tabs`, `upload`). `doc-export` writes bytes.

## Standard flow

1. Confirm task intent and required IDs/filters/paths.
2. Run command with least privilege first when sandboxed execution is available.
3. If command fails and environment is unknown, run `doctor` before retrying.
4. Apply remediation based on exit code and error category.
5. If failure indicates sandbox/runtime config-store or file-read constraints, rerun the same command unsandboxed/escalated.

## Exit-code remediation

- `2` (`validation/config`): fix flags or `--gcloud-bin` path.
  - For upload, verify `--path` is a readable file (not directory/missing path).
  - If error indicates gcloud auth/config store unavailability (for example `gcloud auth config store unavailable`, credential/config-store write/read failures), rerun the same command with escalated/unsandboxed execution.
- `3` (`auth/scope`): re-authenticate and ensure Drive scope.
  - Run: `gcloud auth login --enable-gdrive-access --update-adc`
  - If it persists: `gcloud auth print-access-token`
  - Optional diagnostic: `gcloud auth application-default print-access-token`
- `4` (`network`): retry with larger timeout (for example `--timeout 60s`) and verify connectivity.
- `5` (`api`):
  - Read commands: validate IDs/query/mime/access.
  - Upload: validate destination folder access (`--parent-id`) and shared-drive membership/permissions.

## Safety rules

- Default to read-only operations unless the user explicitly asks for upload.
- Never invent file IDs or document IDs.
- If exporting to disk, use explicit user-provided or clearly scoped output paths.
- For upload, use explicit user-provided local file paths and destination folder IDs.
- If upload path is outside allowed sandbox roots, request escalated execution instead of guessing alternatives.
- Do not claim success without command evidence.

## Response style

- Provide exact command(s) run.
- State invocation mode used (`binary` or `source fallback`).
- Summarize key output fields only.
- Include next remediation step when command fails.
- For escalated retries, report: original command, failing exit code/category, escalated retry command, and final outcome/output path.
