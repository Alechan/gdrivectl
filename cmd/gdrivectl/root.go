package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/Alechan/gdrivectl/internal/app"
	"github.com/Alechan/gdrivectl/internal/fail"
)

func Execute(args []string, stdout, stderr io.Writer) int {
	opts, cmd, cmdArgs, err := parseRootArgs(args)
	if err != nil {
		writeError(stderr, err)
		return fail.ExitCode(err)
	}
	if cmd == "" {
		printUsage(stderr)
		return fail.CodeValidation
	}

	cfg := app.NewConfig(opts.gcloudBin, opts.timeout, opts.json, opts.debug)
	svcs := app.NewServices(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	switch cmd {
	case "doctor":
		return runDoctorCmd(ctx, svcs, cfg, cmdArgs, stdout, stderr)
	case "search":
		return runSearchCmd(ctx, svcs, cfg, cmdArgs, stdout, stderr)
	case "file-meta":
		return runFileMetaCmd(ctx, svcs, cfg, cmdArgs, stdout, stderr)
	case "doc-tabs":
		return runDocTabsCmd(ctx, svcs, cfg, cmdArgs, stdout, stderr)
	case "doc-export":
		return runDocExportCmd(ctx, svcs, cfg, cmdArgs, stdout, stderr)
	case "help", "--help", "-h":
		printUsage(stdout)
		return fail.CodeOK
	default:
		err := fail.NewValidation("unknown command", "use one of: doctor, search, file-meta, doc-tabs, doc-export")
		writeError(stderr, err)
		return fail.ExitCode(err)
	}
}

type rootOptions struct {
	gcloudBin string
	timeout   time.Duration
	json      bool
	debug     bool
}

func parseRootArgs(args []string) (rootOptions, string, []string, error) {
	opts := rootOptions{
		gcloudBin: "gcloud",
		timeout:   20 * time.Second,
	}
	rest := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-h" || a == "--help":
			return opts, "help", nil, nil
		case a == "--json":
			opts.json = true
		case a == "--debug":
			opts.debug = true
		case a == "--gcloud-bin":
			if i+1 >= len(args) {
				return opts, "", nil, fail.NewValidation("missing value for --gcloud-bin", "provide a valid gcloud binary path")
			}
			i++
			opts.gcloudBin = args[i]
		case strings.HasPrefix(a, "--gcloud-bin="):
			opts.gcloudBin = strings.TrimPrefix(a, "--gcloud-bin=")
		case a == "--timeout":
			if i+1 >= len(args) {
				return opts, "", nil, fail.NewValidation("missing value for --timeout", "provide a duration like 20s")
			}
			i++
			d, err := time.ParseDuration(args[i])
			if err != nil {
				return opts, "", nil, fail.NewValidation("invalid --timeout value", "provide a duration like 20s")
			}
			opts.timeout = d
		case strings.HasPrefix(a, "--timeout="):
			d, err := time.ParseDuration(strings.TrimPrefix(a, "--timeout="))
			if err != nil {
				return opts, "", nil, fail.NewValidation("invalid --timeout value", "provide a duration like 20s")
			}
			opts.timeout = d
		default:
			// Keep unknown flags/args for command-level parsing.
			// This lets command flags appear before or after the command token.
			rest = append(rest, a)
		}
	}
	if len(rest) == 0 {
		return opts, "", nil, nil
	}
	return opts, rest[0], rest[1:], nil
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage: gdrivectl [global flags] <command> [flags]")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Commands:")
	fmt.Fprintln(w, "  doctor      Validate gcloud auth and API reachability")
	fmt.Fprintln(w, "  search      Search files in Google Drive")
	fmt.Fprintln(w, "  file-meta   Get metadata for a file")
	fmt.Fprintln(w, "  doc-tabs    List tabs of a Google Doc")
	fmt.Fprintln(w, "  doc-export  Export a Google Doc")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Global flags:")
	fmt.Fprintln(w, "  --gcloud-bin <path>  Path to gcloud binary")
	fmt.Fprintln(w, "  --timeout <duration> Timeout per command (default 20s)")
	fmt.Fprintln(w, "  --json               JSON output")
	fmt.Fprintln(w, "  --debug              Debug logging")
}

func writeError(w io.Writer, err error) {
	var e *fail.Error
	if errors.As(err, &e) {
		if e.Action == "" {
			fmt.Fprintf(w, "Error [%s]: %s\n", e.Category, e.Message)
			return
		}
		fmt.Fprintf(w, "Error [%s]: %s\nAction: %s\n", e.Category, e.Message, e.Action)
		return
	}
	fmt.Fprintf(w, "Error: %s\n", strings.TrimSpace(err.Error()))
}
