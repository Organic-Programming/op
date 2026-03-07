package progress

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/term"
)

type Reporter interface {
	Step(msg string)
	Stepf(format string, args ...any)
	Child() Reporter
}

type Printer struct {
	mu      sync.Mutex
	w       io.Writer
	tty     bool
	silent  bool
	indent  string
	start   time.Time
	now     func() time.Time
	lastLen int
}

func New(w *os.File) *Printer {
	if w == nil {
		return Silence()
	}
	return newPrinter(w, term.IsTerminal(int(w.Fd())), time.Now)
}

func Silence() *Printer {
	return &Printer{
		silent: true,
		start:  time.Now(),
		now:    time.Now,
	}
}

func (p *Printer) Step(msg string) {
	if p == nil || p.silent {
		return
	}
	p.writeLine(p.formatLine(msg), false)
}

func (p *Printer) Stepf(format string, args ...any) {
	p.Step(fmt.Sprintf(format, args...))
}

func (p *Printer) Done(msg string, err error) {
	if p == nil || p.silent {
		return
	}
	mark := "✓"
	if err != nil {
		mark = "✗"
	}
	p.writeLine(p.formatLine(mark+" "+msg), true)
}

func (p *Printer) Child() Reporter {
	if p == nil {
		return Silence()
	}
	return &Printer{
		w:      p.w,
		tty:    p.tty,
		silent: p.silent,
		indent: p.indent + "  ",
		start:  p.start,
		now:    p.now,
	}
}

func (p *Printer) Elapsed() time.Duration {
	if p == nil {
		return 0
	}
	return p.now().Sub(p.start)
}

func FormatTimer(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	total := int(d / time.Second)
	hours := total / 3600
	minutes := (total % 3600) / 60
	seconds := total % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func FormatElapsed(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	seconds := int((d + 500*time.Millisecond) / time.Second)
	if seconds <= 0 {
		seconds = 0
	}
	return fmt.Sprintf("%ds", seconds)
}

func newPrinter(w io.Writer, tty bool, now func() time.Time) *Printer {
	if now == nil {
		now = time.Now
	}
	return &Printer{
		w:     w,
		tty:   tty,
		start: now(),
		now:   now,
	}
}

func (p *Printer) formatLine(msg string) string {
	return fmt.Sprintf("%s%s %s", p.indent, FormatTimer(p.Elapsed()), msg)
}

func (p *Printer) writeLine(line string, final bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.tty && !final {
		padding := ""
		if len(line) < p.lastLen {
			padding = strings.Repeat(" ", p.lastLen-len(line))
		}
		fmt.Fprintf(p.w, "\r%s%s", line, padding)
		p.lastLen = len(line)
		return
	}

	if p.tty && final {
		padding := ""
		if len(line) < p.lastLen {
			padding = strings.Repeat(" ", p.lastLen-len(line))
		}
		fmt.Fprintf(p.w, "\r%s%s\n", line, padding)
		p.lastLen = 0
		return
	}

	fmt.Fprintln(p.w, line)
}
