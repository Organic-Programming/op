package cli

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/organic-programming/grace-op/internal/holons"
	"github.com/organic-programming/grace-op/internal/progress"
	"github.com/organic-programming/grace-op/internal/suggest"
)

type uiOptions struct {
	Quiet bool
}

func commandProgress(format Format, quiet bool) *progress.Printer {
	if quiet || format == FormatJSON {
		return progress.Silence()
	}
	return progress.New(os.Stderr)
}

func extractQuietFlag(args []string) (uiOptions, []string, error) {
	var (
		opts      uiOptions
		remaining []string
	)

	for _, arg := range args {
		switch arg {
		case "--quiet", "-q":
			opts.Quiet = true
		default:
			remaining = append(remaining, arg)
		}
	}
	return opts, remaining, nil
}

func emitSuggestions(w io.Writer, format Format, quiet bool, ctx suggest.Context) {
	if quiet || format == FormatJSON {
		return
	}
	suggest.Print(w, ctx)
}

func humanElapsed(p *progress.Printer) string {
	if p == nil {
		return "0s"
	}
	return progress.FormatElapsed(p.Elapsed())
}

func manifestForSuggestions(ref string) (*holons.LoadedManifest, string) {
	target, err := holons.ResolveTarget(strings.TrimSpace(ref))
	if err != nil || target == nil || target.Manifest == nil {
		return nil, ""
	}
	holon := target.Manifest.BinaryName()
	if holon == "" {
		holon = filepath.Base(target.Dir)
	}
	return target.Manifest, holon
}
