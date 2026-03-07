package progress

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestPrinterNonTTYPrintsSequentialLines(t *testing.T) {
	var buf bytes.Buffer
	base := time.Date(2026, 3, 7, 10, 0, 0, 0, time.UTC)
	now := base
	p := newPrinter(&buf, false, func() time.Time { return now })

	p.Step("checking manifest...")
	now = now.Add(2 * time.Second)
	p.Step("validating prerequisites...")
	now = now.Add(1 * time.Second)
	p.Done("built rob-go in 3s", nil)

	got := buf.String()
	if !strings.Contains(got, "00:00:00 checking manifest...\n") {
		t.Fatalf("output missing first line: %q", got)
	}
	if !strings.Contains(got, "00:00:02 validating prerequisites...\n") {
		t.Fatalf("output missing second line: %q", got)
	}
	if !strings.Contains(got, "00:00:03 ✓ built rob-go in 3s\n") {
		t.Fatalf("output missing done line: %q", got)
	}
}

func TestPrinterTTYUsesCarriageReturn(t *testing.T) {
	var buf bytes.Buffer
	base := time.Date(2026, 3, 7, 10, 0, 0, 0, time.UTC)
	now := base
	p := newPrinter(&buf, true, func() time.Time { return now })

	p.Step("building...")
	now = now.Add(5 * time.Second)
	p.Done("build failed in 5s", errors.New("boom"))

	got := buf.String()
	if !strings.Contains(got, "\r00:00:00 building...") {
		t.Fatalf("tty output missing carriage-return step: %q", got)
	}
	if !strings.Contains(got, "\r00:00:05 ✗ build failed in 5s\n") {
		t.Fatalf("tty output missing final carriage-return line: %q", got)
	}
}

func TestPrinterChildIndentsAndSharesTimer(t *testing.T) {
	var buf bytes.Buffer
	base := time.Date(2026, 3, 7, 10, 0, 0, 0, time.UTC)
	now := base
	p := newPrinter(&buf, false, func() time.Time { return now })
	child, ok := p.Child().(*Printer)
	if !ok {
		t.Fatal("child reporter is not a *Printer")
	}

	now = now.Add(1 * time.Second)
	child.Step("go build -o .op/build/bin/child ./cmd/child")

	got := buf.String()
	if !strings.Contains(got, "  00:00:01 go build -o .op/build/bin/child ./cmd/child\n") {
		t.Fatalf("child output missing indentation: %q", got)
	}
}

func TestSilenceProducesNoOutput(t *testing.T) {
	var buf bytes.Buffer
	p := Silence()
	p.w = &buf
	p.Step("ignored")
	p.Done("ignored", nil)
	if buf.Len() != 0 {
		t.Fatalf("silent printer wrote output: %q", buf.String())
	}
}

func TestFormatHelpers(t *testing.T) {
	if got := FormatTimer(3661 * time.Second); got != "01:01:01" {
		t.Fatalf("FormatTimer = %q, want %q", got, "01:01:01")
	}
	if got := FormatElapsed(3500 * time.Millisecond); got != "4s" {
		t.Fatalf("FormatElapsed = %q, want %q", got, "4s")
	}
}
