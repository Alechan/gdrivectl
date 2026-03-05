# gdrivectl

`gdrivectl` is a spec-driven CLI for reliable Google Drive and Google Docs operations from Codex and terminal workflows.

## Build

```bash
go build ./...
```

## Quickstart

Authenticate with Drive scope first:

```bash
gcloud auth login --enable-gdrive-access --update-adc
```

If needed, use absolute gcloud path:

```bash
GCLOUD=gcloud
```

Doctor check:

```bash
go run ./cmd/gdrivectl doctor --gcloud-bin "$GCLOUD"
```

Search:

```bash
go run ./cmd/gdrivectl search --query "name contains 'RFC'" --page-size 5 --json
```

File metadata:

```bash
go run ./cmd/gdrivectl file-meta --id <FILE_ID> --json
```

Doc tabs:

```bash
go run ./cmd/gdrivectl doc-tabs --id <DOC_ID> --json
```

Doc export:

```bash
go run ./cmd/gdrivectl doc-export --id <DOC_ID> --mime text/plain --out /tmp/doc.txt
```

Optional end-to-end smoke harness:

```bash
GDRIVECTL_FILE_ID=<FILE_ID> GDRIVECTL_DOC_ID=<DOC_ID> scripts/smoke_integration.sh
```

## Commands

- `doctor`: Validate gcloud binary, token retrieval, and Drive/Docs reachability.
- `search`: Query Drive files.
- `file-meta`: Read metadata by file id.
- `doc-tabs`: Read Google Docs tabs metadata.
- `doc-export`: Export a Google Doc to a MIME type.

## Global flags

- `--gcloud-bin <path>`: gcloud binary path.
- `--timeout <duration>`: per-command timeout (default `20s`).
- `--json`: JSON output for structured commands; `doctor` switches from text to JSON.
- `--debug`: enable debug logging.

## Output contract (v0.1)

- `search`, `file-meta`, `doc-tabs`: JSON payload on success.
- `doctor`: text by default, JSON with `--json`.
- `doc-export`: raw bytes to stdout or `--out`; no JSON envelope on success.
- Errors are emitted to `stderr` with category and action hint when available.

## Exit codes

- `0`: success
- `2`: validation/config
- `3`: auth/scope
- `4`: network timeout/reachability
- `5`: API semantic error

## Troubleshooting

### gcloud binary not found

- Check:
  - `command -v gcloud`
- Use `--gcloud-bin <absolute_path>`.

### auth/scope issues

- Check:
  - `gcloud auth list`
- Run `gcloud auth login --enable-gdrive-access --update-adc`.

### timeout/network issues

- Retry with `--timeout 60s`.
- Verify DNS/connectivity.

### Sandbox/config permission issues

- If gcloud cannot write under `~/.config/gcloud`, retry in a shell with home-config access.
- Fallback:
  - `CLOUDSDK_CONFIG=/tmp/gcloud-config gcloud auth list`

### Canonical debug sequence

```bash
go run ./cmd/gdrivectl --help
command -v gcloud
gcloud auth list
go run ./cmd/gdrivectl doctor --json --gcloud-bin "$(command -v gcloud || echo gcloud)"
```

See full guide: `docs/DEBUG.md`.

## SDD-first workflow

1. Start with or update an SDD in `docs/sdd/`.
2. Propose significant design changes in `docs/rfc/`.
3. Record irreversible choices as ADRs in `docs/adr/`.
4. Implement code only after spec acceptance.

## Repository structure

- `docs/sdd/` Software Design Documents (versioned)
- `docs/rfc/` Feature/design proposals before implementation
- `docs/adr/` Architecture Decision Records
- `docs/templates/` Reusable templates for SDD/RFC/ADR
- `cmd/gdrivectl/` CLI entrypoint (Go)
- `internal/` Internal implementation packages
- `scripts/` Local automation scripts (dev/release/docs)

See first SDD iteration: `docs/sdd/0001-gdrivectl-sdd.md`.
