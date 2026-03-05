# v1 Eval Prompts

Use these prompts for quick with-skill vs without-skill evaluation.

1. "Find Drive files with 'RFC' in the name and return first 5 as JSON."
2. "Get metadata for file id <FILE_ID>."
3. "List tabs for Google Doc <DOC_ID>."
4. "Export Google Doc <DOC_ID> as text to /tmp/doc.txt."
5. "I got exit code 3 from file-meta. What should I run next?"
6. "My doc-tabs call times out. Diagnose and retry safely."
7. "Run a preflight check for Drive/Docs readiness."
8. "Use gdrivectl to diagnose why search is failing with code 5."
9. "Given this query, return a command with explicit timeout and JSON output."
10. "I only have a file name, not an ID. Find likely matches then fetch metadata for top candidate."

## Pass criteria

- Correct command selected for intent.
- Correct use of IDs/flags and `--json` where applicable.
- Correct remediation based on exit-code category.
- No write/destructive behavior.
